package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AdminRouteController struct {
	adminController controllers.AdminController
}

func NewAdminRouteController(adminController controllers.AdminController) *AdminRouteController {
	return &AdminRouteController{adminController}
}

func (rc *AdminRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/admin")
	router.GET("/users", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.GetAllUsers)
	router.GET("/users/:id", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.GetUser)
	router.POST("/users/:id/block", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.BlockUser)
	router.POST("/users/:id/unblock", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.UnblockUser)
	router.POST("/categories", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.CreateCategory)
}
