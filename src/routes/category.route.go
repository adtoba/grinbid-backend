package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CategoryRouteController struct {
	categoryController controllers.CategoryController
}

func NewCategoryRouteController(categoryController controllers.CategoryController) *CategoryRouteController {
	return &CategoryRouteController{categoryController}
}

func (rc *CategoryRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/categories")
	router.POST("/", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.categoryController.CreateCategory)
	router.GET("/", middleware.AuthMiddleware(redisClient), rc.categoryController.GetAllCategories)
}
