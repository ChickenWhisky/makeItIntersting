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
	router.GET("/trades/:noOfContracts", GetLastTrades)
	//router.GET("/top/", GetTopOfOrderBook)
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

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "contract": contract})
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

// GetLastTrades returns the current state of the order book
func GetLastTrades(c *gin.Context) {
	_n := c.Param("noOfContracts")
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

// GetTopOfOrderBook gets the top ask and bid level details
//func GetTopOfOrderBook(c *gin.Context) {
//	topOfOrderBook := ob.GetTopOfOrderBook()
//	if len(topOfOrderBook) == 0 {
//		c.JSON(http.StatusOK, gin.H{"message": "No orders in the system"})
//		return
//	} else if len(topOfOrderBook) == 1 {
//		firstLevel := topOfOrderBook[0]
//		if firstLevel.Type {
//			c.JSON(http.StatusOK, gin.H{"message": "Top of the order book", "top_ask": firstLevel})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"message": "Top of the order book", "top_bid": firstLevel})
//		}
//	} else {
//		firstLevel := topOfOrderBook[0]
//		secondLevel := topOfOrderBook[1]
//		if firstLevel.Type {
//			c.JSON(http.StatusOK, gin.H{"message": "Top of the order book", "top_ask": firstLevel, "top_bid": secondLevel})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"message": "Top of the order book", "top_bid": firstLevel, "top_ask": secondLevel})
//		}
//	}
//}
