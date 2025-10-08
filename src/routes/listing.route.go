package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ListingRouteController struct {
	listingController controllers.ListingController
}

func NewListingRouteController(listingController controllers.ListingController) *ListingRouteController {
	return &ListingRouteController{listingController}
}

func (rc *ListingRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/listings")
	router.POST("/", middleware.AuthMiddleware(redisClient), rc.listingController.CreateListing)
	router.GET("/", middleware.AuthMiddleware(redisClient), rc.listingController.GetAllListings)
	router.GET("/:id", middleware.AuthMiddleware(redisClient), rc.listingController.GetListing)
	router.GET("/me", middleware.AuthMiddleware(redisClient), rc.listingController.GetMyListings)
	router.GET("/user/:user_id", middleware.AuthMiddleware(redisClient), rc.listingController.GetAllListingsByUserID)
	router.GET("/category/:category_id", middleware.AuthMiddleware(redisClient), rc.listingController.GetListingByCategory)
}
