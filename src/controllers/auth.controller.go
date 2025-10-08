package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/adtoba/grinbid-backend/src/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB                *gorm.DB
	TokenMaker        *utils.JWTMaker
	SessionController *SessionController
}

func NewAuthController(db *gorm.DB, tokenMaker *utils.JWTMaker, sessionController *SessionController) *AuthController {
	return &AuthController{DB: db, TokenMaker: tokenMaker, SessionController: sessionController}
}

func (ac *AuthController) Login(c *gin.Context) {
	var payload *models.LoginUserRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("invalid credentials: email", nil))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}

	if err := utils.CompareHashAndPassword(payload.Password, user.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("invalid credentials: password", nil))
		return
	}

	accessToken, _, err := ac.TokenMaker.CreateToken(user.ID, user.Email, user.Role, time.Minute*15, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}

	refreshToken, _, err := ac.TokenMaker.CreateToken(user.ID, user.Email, user.Role, time.Hour*24, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}

	res := models.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToUserResponse(),
	}

	c.JSON(http.StatusOK, models.SuccessResponse("login successful", res))
}

func (ac *AuthController) CreateUser(c *gin.Context) {
	var payload *models.CreateUserRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("error creating user", nil))
		return
	}

	user := models.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: hashedPassword,
		Phone:    payload.Phone,
		Location: payload.Location,
		Role:     "user",
	}

	result := ac.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("error creating user", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("user created successfully", user.ToUserResponse()))
}
