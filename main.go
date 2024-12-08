package main

import (
	"github.com/ChickenWhisky/makeItIntersting/docs"
	"github.com/ChickenWhisky/makeItIntersting/internals/handlers"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// Initialize the Gin router
	router := gin.Default()

	// Set up routes from handlers package
	handlers.SetupRoutes(router)

	// use ginSwagger middleware to serve the API docs

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Run the server
	router.Run(":8080")
}
