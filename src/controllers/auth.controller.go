package controllers

import (
	"github.com/adtoba/grinbid-backend/src/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	DB         *gorm.DB
	TokenMaker *utils.JWTMaker
}

func NewAuthController(db *gorm.DB, tokenMaker *utils.JWTMaker) *AuthController {
	return &AuthController{DB: db, TokenMaker: tokenMaker}
}
