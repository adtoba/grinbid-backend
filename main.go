package main

import (
	"fmt"
	"log"

	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/initializers"
	"github.com/adtoba/grinbid-backend/src/migrate"
	"github.com/adtoba/grinbid-backend/src/utils"
)

var (
	AuthController *controllers.AuthController
)

func main() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Println("Config:", config)

	DB := initializers.ConnectDB(&config)
	migrate.Migrate(DB)

	tokenMaker := utils.NewJWTMaker(config.JWT_SECRET)

	AuthController = controllers.NewAuthController(DB, tokenMaker)
}
