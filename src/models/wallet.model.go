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
