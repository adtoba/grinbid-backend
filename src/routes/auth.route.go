package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) *AuthRouteController {
	return &AuthRouteController{authController}
}

func (rc *AuthRouteController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/login", rc.authController.Login)
	router.POST("/register", rc.authController.CreateUser)
}
