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
	router.PUT("/order", ModifyOrder)
	router.DELETE("/order", CancelOrder)
	//router.GET("/orderbook", GetOrderBook)
}

// CreateOrder handles creating a new order
func CreateOrder(c *gin.Context) {
	var contract models.Contract
	if err := c.BindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contract.Timestamp = time.Now().UnixMilli()
	ob.PushContractIntoQueue(contract)

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

// CancelOrder handles canceling an existing order
func CancelOrder(c *gin.Context) {
	var contractForCancellation models.Contract
	if err := c.BindJSON(&contractForCancellation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ob.PushContractIntoQueue(contractForCancellation)

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})

}

// ModifyOrder handles modifying an existing order
func ModifyOrder(c *gin.Context) {
	var contractForModification models.Contract
	if err := c.BindJSON(&contractForModification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ob.PushContractIntoQueue(contractForModification)

	c.JSON(http.StatusOK, gin.H{"message": "Order modified successfully"})
}

// GetOrderBook returns the current state of the order book
//func GetOrderBook(c *gin.Context) {
//	c.JSON(http.StatusOK, ob.GetOrderBook())
//}
