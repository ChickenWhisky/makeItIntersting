package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/events"
	"github.com/ChickenWhisky/makeItIntersting/internals/ledger"
	"github.com/ChickenWhisky/makeItIntersting/pkg/helpers"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// Create a global order book instance
var l = ledger.NewLedger()

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

	// User endpoints
	router.POST("/order", CreateOrder)   // Create a new order
	router.PUT("/order", ModifyOrder)    // Modify an existing order
	router.DELETE("/order", CancelOrder) //	Delete an order
	router.GET("/trades", GetLastTrades) // Get the last n trades for a given event and subevent
	router.GET("/event", GetEvent)       // Get details about a specific event along with its subevents
	router.GET("/events", GetEvents)     // Get list of all current events along with their subevents

	// Admin endpoints

	// Event endpoints
	router.POST("/admin/event", CreateEvent) // Create a new event

	// IMPLEMENT EDITING AN EVENT
	router.PUT("/admin/event")               // Modify an event
	router.DELETE("/admin/event")            // Delete an event

	// SubEvent endpoints
	router.POST("/admin/subevent")   // Create a new subevent
	router.DELETE("/admin/subevent") // Delete a subevent

}

// GetEvents handles getting a list of all current events
func GetEvents(c *gin.Context) {
	events := l.GetEvents()
	c.JSON(http.StatusOK, gin.H{"Events": events})
}

// GetEvent handles getting details about a specific event
func GetEvent(c *gin.Context) {
	
	// Expected JSON format:
	// {
	// 	"event_id": "EventID"
	// }
	
	type TempEvent struct {
		EventID string `json:"event_id"`
	}
	var tempEvent TempEvent
	
	if err := c.BindJSON(&tempEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	
	// Get the event summary
	eventSummary, err := l.GetEvent(tempEvent.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Event": eventSummary})
}

// CreateEvent handles creating a new event
func CreateEvent(c *gin.Context) {
	// Expected JSON format:
	// {
	// 	"eventName": "Event Name",
	// 	"eventInfo": "Event Info",
	// 	"subEvents": [
	// 		"SubEvent1",
	// 		"SubEvent2",
	// 		"SubEvent3",
	// 		"SubEvent4"
	// 	]
	// }

	type TempEvent struct {
		EventName string   `json:"eventName"`
		EventInfo string   `json:"eventInfo"`
		SubEvents []string `json:"subEvents"`
	}
	var tempEvent TempEvent
	if err := c.BindJSON(&tempEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if l.Events[helpers.HashText(tempEvent.EventName)] != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event already exists"})
		return
	}

	// Create a new event
	event, err := events.NewEvents(tempEvent.EventName, tempEvent.SubEvents)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	l.Events[event.GetEventID()] = event
	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully", "event_id": event.GetEventID(), "event_name": event.GetEventName()})
	log.Printf("Event created successfully")
}

// CreateOrder handles creating a new order
func CreateOrder(c *gin.Context) {
	// Expected JSON format:
	// {
	// 	"event_id": "EventID",
	// 	"subevent_id": "SubEventID",
	// 	"user_id": "UserID",
	// 	"order_type": "OrderType",
	// 	"price": 100.0,
	// 	"quantity": 10
	// }

	var Order models.Order
	if err := c.BindJSON(&Order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Order.SetRequestType("add")
	Order.SetTimestamp(time.Now().UnixMilli())

	err := l.SubmitOrder(Order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "Order": Order})
}

// CancelOrder handles canceling an existing order
func CancelOrder(c *gin.Context) {
	// Expected JSON format:
	// {
	// 	"order_id": "OrderID"
	//  "event_id": "EventID",
	// 	"subevent_id": "SubEventID",
	// }
	var OrderForCancellation models.Order
	if err := c.BindJSON(&OrderForCancellation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	OrderForCancellation.SetRequestType("delete")
	err := l.SubmitOrder(OrderForCancellation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})

}

// ModifyOrder handles modifying an existing order
func ModifyOrder(c *gin.Context) {

	// Expected JSON format:
	// {
	// 	"order_id": "OrderID",
	// 	"event_id": "EventID",
	// 	"subevent_id": "SubEventID",
	// 	"price": 100.0,
	// 	"quantity": 10
	// }

	var OrderForModification models.Order

	// Bind the input to the Order struct
	if err := c.BindJSON(&OrderForModification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	OrderForModification.SetRequestType("modify")

	// Submit the modified order
	err := l.SubmitOrder(OrderForModification)

	// Handle errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Order modified successfully"})
}

// For the sake of binding input to a struct for the GetLastTrades function
type GetLastTradesInput struct {
	NoOfOrders int    `json:"no_of_orders"`
	EventID    string `json:"event_id"`
	SubEventID string `json:"subevent_id"`
}

// GetLastTrades returns the current state of the order book
func GetLastTrades(c *gin.Context) {

	// Expected JSON format:
	// {
	// 	"no_of_orders": 10,
	// 	"event_id": "EventID",
	// 	"subevent_id": "SubEventID"
	// }

	// Bind the input to the GetLastTradesInput struct
	var GetLastTradesInput GetLastTradesInput
	if err := c.BindJSON(&GetLastTradesInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Getting the proper order book
	event, eventExists := l.Events[GetLastTradesInput.EventID]
	if !eventExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event doesn't exist"})
		return
	}

	// Getting the proper subevent
	subEvent, subEventExists := event.SubEvents[GetLastTradesInput.SubEventID]
	if !subEventExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SubEvent doesn't exist"})
		return
	}

	// Getting the order book

	ob := subEvent.OrderBook

	n := GetLastTradesInput.NoOfOrders
	lastTradedPrices := ob.GetLastTrades(n)

	if len(lastTradedPrices) != 0 {
		for _, trade := range lastTradedPrices {
			c.JSON(http.StatusOK, trade)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "No trades in the system"})
	}
}
