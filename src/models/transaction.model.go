package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	Amount           float64   `json:"amount" gorm:"not null"`
	ItemName         string    `json:"item_name" gorm:"not null"`
	Status           string    `json:"status" gorm:"not null"`
	Type             string    `json:"type" gorm:"not null"`
	PaymentMethod    string    `json:"payment_method" gorm:"not null"`
	TransactionRef   string    `json:"transaction_ref"`
	ListingID        string    `json:"listing_id"`
	SenderWalletID   string    `json:"sender_wallet_id"`
	ReceiverWalletID string    `json:"receiver_wallet_id"`
	SenderID         string    `json:"sender_id"`
	ReceiverID       string    `json:"receiver_id"`
	Dtt              int       `json:"dtt" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (transaction *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	transaction.ID = uuid.NewString()
	return
}
