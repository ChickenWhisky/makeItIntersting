package main

import (
	"github.com/ChickenWhisky/makeItIntersting/internals/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	

	host := "localhost"
	port := "8080"
	address := host + ":" + port

	router := gin.Default()

	handlers.SetupRoutes(router)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
