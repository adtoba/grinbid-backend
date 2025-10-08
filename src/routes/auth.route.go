package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) *AuthRouteController {
	return &AuthRouteController{authController}
}

func (rc *AuthRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/auth")
	router.POST("/login", rc.authController.Login)
	router.POST("/register", rc.authController.CreateUser)
	router.POST("/refresh-token", rc.authController.RefreshToken)
}
