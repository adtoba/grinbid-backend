package main

import (
	"fmt"
	"log"

	"github.com/adtoba/grinbid-backend/src/initializers"
)

func main() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Println("Config:", config)

	initializers.ConnectDB(&config)
}
