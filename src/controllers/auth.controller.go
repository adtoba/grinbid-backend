package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/adtoba/grinbid-backend/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthController struct {
	DB               *gorm.DB
	TokenMaker       *utils.JWTMaker
	RedisClient      *redis.Client
	WalletController *WalletController
}

func NewAuthController(db *gorm.DB, tokenMaker *utils.JWTMaker, redisClient *redis.Client, walletController *WalletController) *AuthController {
	return &AuthController{DB: db, TokenMaker: tokenMaker, RedisClient: redisClient, WalletController: walletController}
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

	refreshToken, _, err := ac.TokenMaker.CreateToken(user.ID, user.Email, user.Role, time.Hour*168, true)
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

	// Generate unique username
	username, err := utils.GenerateUsername(ac.DB, payload.FullName, payload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("error generating username", nil))
		return
	}

	user := models.User{
		Username: username,
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

	// Create a wallet for the user
	ac.WalletController.CreateWallet(c, user.ID)

	c.JSON(http.StatusCreated, models.SuccessResponse("user created successfully", user.ToUserResponse()))
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	// Bind the request body
	var payload *models.RenewAccessTokenRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	// Check if the refresh token is in the Redis client
	userID, err := ac.RedisClient.Get(c, payload.RefreshToken).Result()
	if err == redis.Nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("invalid or expired refresh token", nil))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}

	// Verify the refresh token
	userClaims, err := ac.TokenMaker.VerifyToken(payload.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("invalid or expired refresh token", nil))
		return
	}

	// Check if the user ID matches the refresh token
	if userClaims.ID != userID {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("invalid or expired refresh token", nil))
		return
	}

	// Revoke the old refresh token
	err = ac.RedisClient.Del(c, payload.RefreshToken).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to revoke old refresh token", nil))
		return
	}

	// Blacklist the old refresh token
	ttl := time.Until(userClaims.RegisteredClaims.ExpiresAt.Time)
	if ttl > 0 {
		err = ac.RedisClient.Set(c, "blacklist:"+payload.RefreshToken, "revoked", ttl).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to blacklist old refresh token", nil))
			return
		}
	}

	// Create a new access token
	newAccessToken, _, err := ac.TokenMaker.CreateToken(userID, userClaims.Email, userClaims.Role, time.Minute*15, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create new access token", nil))
		return
	}

	// Create a new refresh token
	newRefreshToken, _, err := ac.TokenMaker.CreateToken(userID, userClaims.Email, userClaims.Role, time.Hour*168, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create new refresh token", nil))
		return
	}

	// Return the new tokens
	res := models.RenewAccessTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Token renewed successfully", res))

}
