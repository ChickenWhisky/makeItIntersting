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
	router.GET("/trades", GetLastTrades)
	router.GET("/top", GetTopOfOrderBook)
}

// CreateOrder handles creating a new order
func CreateOrder(c *gin.Context) {
	var contract models.Contract
	if err := c.BindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contract.Timestamp = time.Now().UnixMilli()
	ob.AddContract(contract)

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "contract": contract})
}

// CancelOrder handles canceling an existing order
func CancelOrder(c *gin.Context) {
	var contractForCancellation models.Contract
	if err := c.BindJSON(&contractForCancellation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ob.CancelContract(contractForCancellation)

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})

}

// ModifyOrder handles modifying an existing order
func ModifyOrder(c *gin.Context) {
	var contractForModification models.Contract
	if err := c.BindJSON(&contractForModification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ob.ModifyContract(contractForModification)

	c.JSON(http.StatusOK, gin.H{"message": "Order modified successfully"})
}

// GetLastTrades returns the current state of the order book
func GetLastTrades(c *gin.Context) {
	lastTradedPrices := ob.GetLastTrades()
	for _, trade := range *lastTradedPrices {
		c.JSON(http.StatusOK, trade)
	}
}

// GetTopOfOrderBook gets the top ask and bid level details
func GetTopOfOrderBook(c *gin.Context) {
	topOfOrderBook := ob.GetTopOfOrderBook()
	if len(topOfOrderBook) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No orders in the system"})
		return
	} else {
		firstLevel := topOfOrderBook[0]
		secondLevel := topOfOrderBook[1]
		if firstLevel.Type {
			c.JSON(http.StatusOK, gin.H{"message": "Top of the order book", "top_ask": firstLevel, "top_bid": secondLevel})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Top of the order book", "top_bid": firstLevel, "top_ask": secondLevel})
		}
	}
}
