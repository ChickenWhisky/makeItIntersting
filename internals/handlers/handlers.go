package handlers

import (
	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
	"github.com/ChickenWhisky/makeItIntersting/pkg/helpers"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a global order book instance
var ob = orderbook.NewOrderBook()

// SetupRoutes configures the routes for the Gin router
func SetUpCors(router *gin.Engine, web_url string) {
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", web_url)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(204)
	})
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{web_url}, // Explicitly allow your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,

		AllowOriginFunc: func(origin string) bool {
			return origin == web_url
		},
		MaxAge: 12 * time.Hour,
	}))
}
func SetupRoutes(router *gin.Engine) {
	router.POST("/order", CreateOrder)
	router.PUT("/order", ModifyOrder)
	router.DELETE("/order", CancelOrder)
	router.GET("/trades/:noOfOrders", GetLastTrades)
	//router.GET("/top/", GetTopOfOrderBook)
}

// CreateOrder handles creating a new order
func CreateOrder(c *gin.Context) {
	var Order models.Order
	if err := c.BindJSON(&Order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Order.SetRequestType("add")
	Order.SetTimestamp(time.Now().UnixMilli()) 
	ob.PushOrderIntoQueue(Order)

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "Order": Order})
}

// CancelOrder handles canceling an existing order
func CancelOrder(c *gin.Context) {
	var OrderForCancellation models.Order
	if err := c.BindJSON(&OrderForCancellation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	OrderForCancellation.SetRequestType("delete")
	ob.PushOrderIntoQueue(OrderForCancellation)

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})

}

// ModifyOrder handles modifying an existing order
func ModifyOrder(c *gin.Context) {
	var OrderForModification models.Order
	if err := c.BindJSON(&OrderForModification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	OrderForModification.SetRequestType("modify")
	ob.PushOrderIntoQueue(OrderForModification)

	c.JSON(http.StatusOK, gin.H{"message": "Order modified successfully"})
}

// GetLastTrades returns the current state of the order book
func GetLastTrades(c *gin.Context) {
	_n := c.Param("noOfOrders")
	n := helpers.ConvertStringToInt(_n)
	lastTradedPrices := ob.GetLastTrades(n)
	if len(lastTradedPrices) != 0 {
		for _, trade := range lastTradedPrices {
			c.JSON(http.StatusOK, trade)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "No trades in the system"})
	}
}

