package controllers

import (
	"fmt"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SessionController struct {
	DB *gorm.DB
}

func NewSessionController(db *gorm.DB) *SessionController {
	return &SessionController{DB: db}
}

func (sc *SessionController) CreateSession(c *gin.Context, s models.Session) (*models.Session, error) {
	result := sc.DB.Create(&s)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to create session: %w", result.Error)
	}

	return &s, nil
}

func (sc *SessionController) GetSession(c *gin.Context, id string) (*models.Session, error) {
	var s models.Session

	result := sc.DB.First(&s, "id = ?", id)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get session: %w", result.Error)
	}
	return &s, nil
}

func (sc *SessionController) RevokeSession(c *gin.Context, id string) error {
	result := sc.DB.Model(&models.Session{}).Where("id = ?", id).Update("is_revoked", true)

	if result.Error != nil {
		return fmt.Errorf("failed to revoke session: %w", result.Error)
	}
	return nil
}

func (sc *SessionController) DeleteSession(c *gin.Context, id string) error {
	result := sc.DB.Delete(&models.Session{}).Where("id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete session: %w", result.Error)
	}
	return nil
}
