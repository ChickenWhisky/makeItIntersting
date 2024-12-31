package main

import (
	"github.com/ChickenWhisky/makeItIntersting/docs"
	"github.com/ChickenWhisky/makeItIntersting/internals/handlers"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
)

func main() {

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Get host and port from environment variables
	web_url := os.Getenv("WEB_URL")

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	address := host + ":" + port

	// Initialize the Gin router
	router := gin.Default()
	handlers.SetUpCors(router, web_url)
	handlers.SetupRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(address)
}
