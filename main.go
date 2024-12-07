package main

import (
	"github.com/ChickenWhisky/makeItIntersting/internals/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Set up routes from handlers package
	handlers.SetupRoutes(router)

	// Run the server
	router.Run(":8080")
}
