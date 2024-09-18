package handlers

import (
	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a global order book instance
var ob = orderbook.NewOrderBook()

// SetupRoutes configures the routes for the Gin router
func SetupRoutes(router *gin.Engine) {
	router.POST("/order", CreateOrder)
	router.DELETE("/order", CancelOrder)
	router.GET("/orderbook", GetOrderBook)
}

// CreateOrder handles creating a new order
func CreateOrder(c *gin.Context) {
	var contract models.Contract
	if err := c.BindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contract.Timestamp = time.Now().Unix()
	ob.AddContract(contract)
	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

// CancelOrder handles canceling an existing order
func CancelOrder(c *gin.Context) {
	var data struct {
		UserID    string  `json:"user_id"`
		OrderType string  `json:"order_type"`
		Price     float64 `json:"price"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ob.CancelContract(data.UserID, data.OrderType, data.Price)
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// GetOrderBook returns the current state of the order book
func GetOrderBook(c *gin.Context) {
	c.JSON(http.StatusOK, ob.GetOrderBook())
}
