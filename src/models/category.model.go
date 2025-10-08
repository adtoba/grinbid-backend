package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Icon      string    `json:"icon" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (category *Category) BeforeCreate(tx *gorm.DB) (err error) {
	category.ID = uuid.NewString()
	return
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon" binding:"required"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
