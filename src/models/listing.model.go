package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Listing struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Price       int       `json:"price" gorm:"not null"`
	Image       string    `json:"image" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null"`
	Location    string    `json:"location" gorm:"not null"`
	Condition   string    `json:"condition" gorm:"not null"`
	CategoryID  string    `json:"category_id" gorm:"not null"`
	UserID      string    `json:"user_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (listing *Listing) BeforeCreate(tx *gorm.DB) (err error) {
	listing.ID = uuid.NewString()
	return
}

type ListingResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Image       string    `json:"image"`
	Status      string    `json:"status"`
	Location    string    `json:"location"`
	Condition   string    `json:"condition"`
	CategoryID  string    `json:"category_id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
