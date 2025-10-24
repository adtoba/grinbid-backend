package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	AvailableBalance float64   `json:"available_balance" gorm:"not null"`
	LockedBalance    float64   `json:"locked_balance" gorm:"not null"`
	UserID           string    `json:"user_id" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (wallet *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	wallet.ID = uuid.NewString()
	return
}

type PurchaseRequest struct {
	Amount     float64 `json:"amount" binding:"required"`
	SenderID   string  `json:"sender_id" binding:"required"`
	ReceiverID string  `json:"receiver_id" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	ListingID  string  `json:"listing_id" binding:"required"`
	Dtt        int     `json:"dtt" binding:"required" gorm:"default:0"`
}

type WithdrawFromWalletRequest struct {
	Amount        float64 `json:"amount" binding:"required"`
	AccountNumber string  `json:"account_number" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	BankName      string  `json:"bank_name" binding:"required"`
	AccountName   string  `json:"account_name" binding:"required"`
	Ref           string  `json:"ref" binding:"required"`
}

type TopUpWalletRequest struct {
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Ref           string  `json:"ref" binding:"required"`
}

type WalletResponse struct {
	ID               string    `json:"id"`
	AvailableBalance float64   `json:"available_balance"`
	LockedBalance    float64   `json:"locked_balance"`
	UserID           string    `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type InitializeTransactionRequest struct {
	Amount     float64 `json:"amount" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	ItemName   string  `json:"item_name"`
	ReceiverID string  `json:"receiver_id"`
	ListingID  string  `json:"listing_id"`
	Dtt        int     `json:"dtt" gorm:"default:0"`
}
