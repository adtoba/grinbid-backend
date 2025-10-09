// Transaction Types:
// 1. purchase
// 2. topup
// 3. withdraw
// 4. refund

package controllers

import (
	"net/http"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletController struct {
	DB *gorm.DB
}

func NewWalletController(db *gorm.DB) *WalletController {
	return &WalletController{DB: db}
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

// GetWalletTransactions gets all transactions for a wallet

// WithdrawFromWallet withdraws from a wallet

// Purchase purchases a listing
func (wc *WalletController) Purchase(c *gin.Context) {
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

func (wc *WalletController) TopUp(c *gin.Context) {
	var payload models.TopUpWalletRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	// Verify the payment via payment provider

	// Get the wallet
	var wallet models.Wallet
	userID := c.MustGet("user_id").(string)
	result := wc.DB.First(&wallet, "user_id = ?", userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("wallet not found", nil))
		return
	}
	wallet.AvailableBalance += payload.Amount

	var transaction models.Transaction
	transaction.Amount = payload.Amount
	transaction.Type = "topup"
	transaction.TransactionRef = payload.Ref
	transaction.SenderWalletID = wallet.ID
	transaction.SenderID = userID
	transaction.Status = "success"

	result = wc.DB.Save(&wallet)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("error saving wallet", nil))
		return
	}

	result = wc.DB.Create(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("error saving transaction", nil))
		return
	}

	// Send transaction notification to user via push notification

	c.JSON(http.StatusOK, models.SuccessResponse("Topup successful", nil))
}

func (wc *WalletController) Withdraw(c *gin.Context) {
	var payload models.WithdrawFromWalletRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	var wallet models.Wallet
	userID := c.MustGet("user_id").(string)
	result := wc.DB.First(&wallet, "user_id = ?", userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("wallet not found", nil))
		return
	}

	if wallet.AvailableBalance < payload.Amount {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Insufficient balance", nil))
		return
	}

	wallet.AvailableBalance -= payload.Amount

	// Send transaction notification to user via push notification

	// Send money to user's account via payment provider

	var transaction models.Transaction
	transaction.Amount = payload.Amount
	transaction.Type = "withdraw"
	transaction.PaymentMethod = payload.PaymentMethod
	transaction.TransactionRef = uuid.NewString()
	transaction.SenderWalletID = wallet.ID
	transaction.SenderID = userID
	transaction.Status = "success"
	transaction.TransactionRef = payload.Ref

	wc.DB.Save(&wallet)
	wc.DB.Create(&transaction)

	c.JSON(http.StatusOK, models.SuccessResponse("Withdrawal successful", nil))
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
