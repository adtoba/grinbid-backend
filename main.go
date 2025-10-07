package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/initializers"
	"github.com/adtoba/grinbid-backend/src/migrate"
	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/adtoba/grinbid-backend/src/routes"
	"github.com/adtoba/grinbid-backend/src/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      *controllers.AuthController
	AuthRouteController *routes.AuthRouteController
	SessionController   *controllers.SessionController
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Println("Config:", config)

	DB := initializers.ConnectDB(&config)
	migrate.Migrate(DB)

	tokenMaker := utils.NewJWTMaker(config.JWT_SECRET)
	SessionController = controllers.NewSessionController(DB)
	AuthController = controllers.NewAuthController(DB, tokenMaker, SessionController)
	AuthRouteController = routes.NewAuthRouteController(*AuthController)
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

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.SuccessResponse("server is running", nil))
	})

	v1 := router.Group("/api/v1")
	{
		AuthRouteController.RegisterRoutes(v1)
	}

	log.Fatal((server.Run(":" + config.ServerPort)))
}
