package models

import "time"

type Session struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
