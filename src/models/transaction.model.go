package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Amount         float64   `json:"amount" gorm:"not null"`
	Status         string    `json:"status" gorm:"not null"`
	Type           string    `json:"type" gorm:"not null"`
	TransactionRef string    `json:"transaction_ref" gorm:"not null"`
	ListingID      string    `json:"listing_id" gorm:"not null"`
	WalletID       string    `json:"wallet_id" gorm:"not null"`
	SenderID       string    `json:"sender_id" gorm:"not null"`
	ReceiverID     string    `json:"receiver_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (transaction *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	transaction.ID = uuid.NewString()
	return
}
