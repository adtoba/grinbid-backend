package main

import (
	"fmt"
	"log"

	"github.com/adtoba/grinbid-backend/src/initializers"
	"github.com/adtoba/grinbid-backend/src/migrate"
)

func main() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Println("Config:", config)

	DB := initializers.ConnectDB(&config)
	migrate.Migrate(DB)
}
