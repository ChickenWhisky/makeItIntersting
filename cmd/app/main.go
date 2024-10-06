package main

import (
	"github.com/ChickenWhisky/makeItIntersting/internals/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Set up routes from handlers package
	handlers.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
