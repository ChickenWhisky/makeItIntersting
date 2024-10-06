package main

import (
	"github.com/ChickenWhisky/makeItIntersting/internals/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get host and port from environment variables
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	address := host + ":" + port

	// Initialize the Gin router
	router := gin.Default()

	// Set up routes from handlers package
	handlers.SetupRoutes(router)

	// Run the server
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
