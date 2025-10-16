package services

import (
	"errors"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/adtoba/grinbid-backend/src/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WalletService struct {
	DB *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{DB: db}
}

func (ws *WalletService) TopupWallet(c *gin.Context, transaction models.Transaction) error {
	userID := transaction.SenderID

	var wallet models.Wallet
	ws.DB.First(&wallet, "user_id = ?", userID)
	if wallet.ID == "" {
		return errors.New("wallet not found")
	}

	wallet.AvailableBalance += transaction.Amount
	result := ws.DB.Save(&wallet)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ws *WalletService) PurchaseFromWallet(c *gin.Context, transaction models.Transaction) error {
	userID := c.MustGet("user_id").(string)

	var senderWallet models.Wallet
	ws.DB.First(&senderWallet, "user_id = ?", userID)
	if senderWallet.ID == "" {
		return errors.New("wallet not found")
	}

	var receiverWallet models.Wallet
	ws.DB.First(&receiverWallet, "user_id = ?", transaction.ReceiverID)
	if receiverWallet.ID == "" {
		return errors.New("receiver wallet not found")
	}

	if senderWallet.AvailableBalance < transaction.Amount {
		return errors.New("insufficient balance")
	}

	senderWallet.AvailableBalance -= transaction.Amount

	if transaction.Dtt > 0 {
		receiverWallet.LockedBalance += transaction.Amount
	} else {
		receiverWallet.AvailableBalance += transaction.Amount
	}

	ws.DB.Save(&senderWallet)
	ws.DB.Save(&receiverWallet)

	return nil
}

func (ws *WalletService) PurchaseFromPaystack(c *gin.Context, transaction models.Transaction) error {
	var receiverWallet models.Wallet
	ws.DB.First(&receiverWallet, "user_id = ?", transaction.ReceiverID)
	if receiverWallet.ID == "" {
		return errors.New("receiver wallet not found")
	}

	if transaction.Dtt > 0 {
		receiverWallet.LockedBalance += transaction.Amount
	} else {
		receiverWallet.AvailableBalance += transaction.Amount
	}

	ws.DB.Save(&receiverWallet)

	return nil
}

func (ws *WalletService) WithdrawFromWallet(c *gin.Context, transaction models.Transaction) error {
	userID := c.MustGet("user_id").(string)

	var wallet models.Wallet
	ws.DB.First(&wallet, "user_id = ?", userID)
	if wallet.ID == "" {
		return errors.New("wallet not found")
	}

	convertedAmount := utils.FromKobo(transaction.Amount)

	wallet.AvailableBalance -= convertedAmount
	ws.DB.Save(&wallet)

	return nil
}
