package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	FullName        string    `json:"full_name" gorm:"not null"`
	Email           string    `json:"email" gorm:"not null"`
	Password        string    `json:"password" gorm:"not null"`
	Phone           string    `json:"phone" gorm:"not null"`
	Location        string    `json:"address" gorm:"not null"`
	Role            string    `json:"role" gorm:"not null"`
	Nin             string    `json:"nin" gorm:"not null"`
	IsVerified      bool      `json:"is_verified" gorm:"default:false"`
	IsBlocked       bool      `json:"is_blocked" gorm:"default:false"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"default:false"`
	IsPhoneVerified bool      `json:"is_phone_verified" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return
}

type UserResponse struct {
	ID              string    `json:"id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Location        string    `json:"location"`
	Role            string    `json:"role"`
	Nin             string    `json:"nin"`
	IsVerified      bool      `json:"is_verified"`
	IsBlocked       bool      `json:"is_blocked"`
	IsEmailVerified bool      `json:"is_email_verified"`
	IsPhoneVerified bool      `json:"is_phone_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (user *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:              user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone,
		Location:        user.Location,
		Role:            user.Role,
		Nin:             user.Nin,
		IsVerified:      user.IsVerified,
		IsBlocked:       user.IsBlocked,
		IsEmailVerified: user.IsEmailVerified,
		IsPhoneVerified: user.IsPhoneVerified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
