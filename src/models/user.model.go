package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	Username        string    `json:"username" gorm:"uniqueIndex;not null"`
	FullName        string    `json:"full_name" gorm:"not null"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Password        string    `json:"password" gorm:"not null"`
	Phone           string    `json:"phone" gorm:"not null"`
	Location        string    `json:"address" gorm:"not null"`
	Role            string    `json:"role" gorm:"not null"`
	Nin             string    `json:"nin"`
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
	Username        string    `json:"username"`
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

type SimpleUserResponse struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Location        string    `json:"location"`
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
		Username:        user.Username,
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

func (user *User) ToSimpleUserResponse() SimpleUserResponse {
	return SimpleUserResponse{
		ID:              user.ID,
		Username:        user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone,
		Location:        user.Location,
		IsVerified:      user.IsVerified,
		IsBlocked:       user.IsBlocked,
		IsEmailVerified: user.IsEmailVerified,
		IsPhoneVerified: user.IsPhoneVerified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone" binding:"required"`
	Location string `json:"location"`
}

type LoginUserResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
