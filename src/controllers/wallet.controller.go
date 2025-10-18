// Transaction Types:
// 1. purchase
// 2. topup
// 3. withdraw
// 4. refund

package controllers

import (
	"net/http"
	"strconv"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/adtoba/grinbid-backend/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletController struct {
	DB              *gorm.DB
	PaystackService *services.PaystackService
}

func NewWalletController(db *gorm.DB, paystackService *services.PaystackService) *WalletController {
	return &WalletController{DB: db, PaystackService: paystackService}
}

// CreateWallet creates a wallet for a user
func (wc *WalletController) CreateWallet(c *gin.Context, id string) {
	var wallet models.Wallet
	userID, exists := c.Get("user_id")
	if !exists {
		wallet.UserID = id
	} else {
		wallet.UserID = userID.(string)
	}

	wallet.AvailableBalance = 0
	wallet.LockedBalance = 0
	result := wc.DB.Create(&wallet)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("error creating wallet", nil))
		return
	}
}

// GetWallet gets a wallet for a user
func (wc *WalletController) GetWallet(c *gin.Context) {
	var wallet models.Wallet
	userID := c.MustGet("user_id").(string)
	result := wc.DB.First(&wallet, "user_id = ?", userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("wallet not found", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("wallet fetched successfully", wallet))
}

// Purchase purchases a listing
func (wc *WalletController) PurchaseFromWallet(c *gin.Context) {
	var payload models.PurchaseRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	var senderWallet models.Wallet
	wc.DB.First(&senderWallet, "user_id = ?", payload.SenderID)
	if senderWallet.ID == "" {
		c.JSON(http.StatusNotFound, models.ErrorResponse("sender wallet not found", nil))
		return
	}

	var receiverWallet models.Wallet
	wc.DB.First(&receiverWallet, "user_id = ?", payload.ReceiverID)
	if receiverWallet.ID == "" {
		c.JSON(http.StatusNotFound, models.ErrorResponse("receiver wallet not found", nil))
		return
	}

	if senderWallet.AvailableBalance < payload.Amount {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Insufficient balance", nil))
		return
	}
	senderWallet.AvailableBalance -= payload.Amount
	if payload.Dtt > 0 {
		receiverWallet.LockedBalance += payload.Amount
	} else {
		receiverWallet.AvailableBalance += payload.Amount
	}

	var transaction models.Transaction
	transaction.Amount = payload.Amount
	transaction.Type = "purchase"
	transaction.ListingID = payload.ListingID
	transaction.SenderID = payload.SenderID
	transaction.ReceiverID = payload.ReceiverID
	transaction.SenderWalletID = senderWallet.ID
	transaction.ReceiverWalletID = receiverWallet.ID
	transaction.TransactionRef = uuid.NewString()
	transaction.Dtt = payload.Dtt

	if payload.Dtt > 0 {
		transaction.Status = "in_progress"
	} else {
		transaction.Status = "success"
	}

	wc.DB.Save(&senderWallet)
	wc.DB.Save(&receiverWallet)
	wc.DB.Create(&transaction)

	// Send transaction notification to sender and receiver via chat and push notification

	c.JSON(http.StatusOK, models.SuccessResponse("Wallet charged successfully", nil))
}

func (wc *WalletController) InitializeTransaction(c *gin.Context) {
	var payload models.InitializeTransactionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	userEmail := c.MustGet("user_email").(string)
	userID := c.MustGet("user_id").(string)

	var senderWallet models.Wallet
	wc.DB.First(&senderWallet, "user_id = ?", userID)
	if senderWallet.ID == "" {
		c.JSON(http.StatusNotFound, models.ErrorResponse("sender wallet not found", nil))
		return
	}

	var receiverWallet models.Wallet

	if payload.ReceiverID == "" {
		payload.ReceiverID = userID
	} else {
		wc.DB.First(&receiverWallet, "user_id = ?", payload.ReceiverID)
		if receiverWallet.ID == "" {
			c.JSON(http.StatusNotFound, models.ErrorResponse("receiver wallet not found", nil))
			return
		}
	}

	var transaction models.Transaction
	transaction.Amount = payload.Amount
	transaction.Type = payload.Type
	transaction.Status = "initiated"
	transaction.TransactionRef = uuid.NewString()
	transaction.SenderID = userID
	transaction.ReceiverID = payload.ReceiverID
	transaction.ListingID = payload.ListingID
	transaction.Dtt = payload.Dtt
	transaction.SenderWalletID = senderWallet.ID
	transaction.ReceiverWalletID = receiverWallet.ID

	response, err := wc.PaystackService.InitializeTransaction(strconv.FormatFloat(payload.Amount, 'f', -1, 64), userEmail, transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("error initializing transaction", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Transaction initialized successfully", response))
}

func (wc *WalletController) InitiateWithdrawal(c *gin.Context) {

}

func (wc *WalletController) GetWalletTransactions(c *gin.Context) {
	var transactions []models.Transaction
	userID := c.MustGet("user_id").(string)
	result := wc.DB.Find(&transactions, "sender_id = ? OR receiver_id = ?", userID, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("transactions not found", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Transactions fetched successfully", transactions))
}

func (wc *WalletController) GetAllWalletTransactions(c *gin.Context) {
	var transactions []models.Transaction
	result := wc.DB.Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("transactions not found", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Transactions fetched successfully", transactions))
}

func (wc *WalletController) GetWalletTransactionById(c *gin.Context) {
	var transaction models.Transaction
	id := c.Param("id")
	result := wc.DB.First(&transaction, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("transaction not found", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Transaction fetched successfully", transaction))
}

func (wc *WalletController) GetWalletTransactionsByUserId(c *gin.Context) {
	var transactions []models.Transaction
	id := c.Param("id")
	result := wc.DB.Find(&transactions, "sender_id = ? OR receiver_id = ?", id, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("transactions not found", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Transactions fetched successfully", transactions))
}
