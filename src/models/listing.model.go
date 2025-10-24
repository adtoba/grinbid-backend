package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Listing struct {
	ID          string             `json:"id" gorm:"primaryKey"`
	Title       string             `json:"title" gorm:"not null"`
	Description string             `json:"description" gorm:"not null"`
	Price       float64            `json:"price" gorm:"not null"`
	Image       string             `json:"image" gorm:"not null"`
	Status      string             `json:"status" gorm:"not null"`
	Location    string             `json:"location" gorm:"not null"`
	Condition   string             `json:"condition" gorm:"not null"`
	CategoryID  string             `json:"category_id" gorm:"not null"`
	UserID      string             `json:"user_id" gorm:"not null"`
	User        SimpleUserResponse `json:"user"`
	CreatedAt   time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
}

func (listing *Listing) BeforeCreate(tx *gorm.DB) (err error) {
	listing.ID = uuid.NewString()
	return
}

type CreateListingRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Image       string  `json:"image" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	Condition   string  `json:"condition" binding:"required"`
	CategoryID  string  `json:"category_id" binding:"required"`
}

type ListingResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       float64            `json:"price"`
	Image       string             `json:"image"`
	Status      string             `json:"status"`
	Location    string             `json:"location"`
	Condition   string             `json:"condition"`
	CategoryID  string             `json:"category_id"`
	UserID      string             `json:"user_id"`
	User        SimpleUserResponse `json:"user"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func (listing *Listing) CreateResponse() ListingResponse {
	return ListingResponse{
		ID:          listing.ID,
		Title:       listing.Title,
		Description: listing.Description,
		Price:       listing.Price,
		Image:       listing.Image,
		Status:      listing.Status,
		Location:    listing.Location,
		Condition:   listing.Condition,
		CategoryID:  listing.CategoryID,
		UserID:      listing.UserID,
		CreatedAt:   listing.CreatedAt,
		UpdatedAt:   listing.UpdatedAt,
	}
}
