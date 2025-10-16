package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/initializers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/adtoba/grinbid-backend/src/migrate"
	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/adtoba/grinbid-backend/src/routes"
	"github.com/adtoba/grinbid-backend/src/services"
	"github.com/adtoba/grinbid-backend/src/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	server      *gin.Engine
	RedisClient *redis.Client

	WalletController      *controllers.WalletController
	WalletRouteController *routes.WalletRouteController

	AuthController      *controllers.AuthController
	AuthRouteController *routes.AuthRouteController

	ListingController      *controllers.ListingController
	ListingRouteController *routes.ListingRouteController

	AdminController      *controllers.AdminController
	AdminRouteController *routes.AdminRouteController

	CategoryController      *controllers.CategoryController
	CategoryRouteController *routes.CategoryRouteController

	WebhookController      *controllers.WebhooksController
	WebhookRouteController *routes.WebhookRouteController

	PaystackService *services.PaystackService
	WalletService   *services.WalletService

	ctx = context.Background()
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Println("Config:", config)

	DB := initializers.ConnectDB(&config)
	migrate.Migrate(DB)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	tokenMaker := utils.NewJWTMaker(config.JWT_SECRET, RedisClient)

	PaystackService = services.NewPaystackService(config.PaystackSecretKey)
	WalletService = services.NewWalletService(DB)

	WalletController = controllers.NewWalletController(DB, PaystackService)
	WalletRouteController = routes.NewWalletRouteController(*WalletController)

	AuthController = controllers.NewAuthController(DB, tokenMaker, RedisClient, WalletController)
	AuthRouteController = routes.NewAuthRouteController(*AuthController)

	ListingController = controllers.NewListingController(DB)
	ListingRouteController = routes.NewListingRouteController(*ListingController)

	AdminController = controllers.NewAdminController(DB)
	AdminRouteController = routes.NewAdminRouteController(*AdminController, *WalletController)

	CategoryController = controllers.NewCategoryController(DB)
	CategoryRouteController = routes.NewCategoryRouteController(*CategoryController)

	WebhookController = controllers.NewWebhooksController(DB, WalletService)
	WebhookRouteController = routes.NewWebhookRouteController(*WebhookController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	router := server

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // In production, replace with specific origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	router.GET("/health", middleware.AuthMiddleware(RedisClient), func(c *gin.Context) {
		c.JSON(http.StatusOK, models.SuccessResponse("server is running", nil))
	})

	v1 := router.Group("/api/v1")
	{
		AuthRouteController.RegisterRoutes(v1, RedisClient)
		ListingRouteController.RegisterRoutes(v1, RedisClient)
		AdminRouteController.RegisterRoutes(v1, RedisClient)
		CategoryRouteController.RegisterRoutes(v1, RedisClient)
		WalletRouteController.RegisterRoutes(v1, RedisClient)
		WebhookRouteController.RegisterRoutes(v1, RedisClient)
	}

	log.Fatal((server.Run(":" + config.ServerPort)))
}
