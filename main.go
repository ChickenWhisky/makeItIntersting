package main

import (
	"net/http"
	"time"
	"sort"
	"sync"
	"github.com/gin-gonic/gin"
)

// Contract struct - represents a buy/sell order
type Contract struct {
	UserID    string  `json:"user_id"`
	OrderType string  `json:"order_type"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Timestamp int64   `json:"timestamp"`
}

// OrderBook struct - stores the order book data and manages orders
type OrderBook struct {
	Asks             map[float64]int
	Bids             map[float64]int
	UserOrders       map[string][]Contract
	LastMatchedPrice *float64
	mu               sync.Mutex // to handle concurrent access
}

// NewOrderBook creates a new empty order book
func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks:       make(map[float64]int),
		Bids:       make(map[float64]int),
		UserOrders: make(map[string][]Contract),
	}
}

// AddContract adds a new contract (order) to the order book and attempts to match orders
func (ob *OrderBook) AddContract(contract Contract) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	switch contract.OrderType {
	case "sell":
		ob.Asks[contract.Price] += contract.Quantity
	case "buy":
		ob.Bids[contract.Price] += contract.Quantity
	}

	// Store the user's orders
	ob.UserOrders[contract.UserID] = append(ob.UserOrders[contract.UserID], contract)

	// Attempt to match orders
	ob.matchOrders()
}

// matchOrders matches highest bid with the lowest ask
func (ob *OrderBook) matchOrders() {
	for len(ob.Asks) > 0 && len(ob.Bids) > 0 {
		lowestAskPrice := ob.getLowestAskPrice()
		highestBidPrice := ob.getHighestBidPrice()

		if highestBidPrice < lowestAskPrice {
			break
		}

		matchedQuantity := min(ob.Asks[lowestAskPrice], ob.Bids[highestBidPrice])

		ob.LastMatchedPrice = &lowestAskPrice
		ob.Asks[lowestAskPrice] -= matchedQuantity
		ob.Bids[highestBidPrice] -= matchedQuantity

		if ob.Asks[lowestAskPrice] == 0 {
			delete(ob.Asks, lowestAskPrice)
		}
		if ob.Bids[highestBidPrice] == 0 {
			delete(ob.Bids, highestBidPrice)
		}
	}
}

// CancelContract cancels a specific user's contract
func (ob *OrderBook) CancelContract(userID string, orderType string, price float64) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if contracts, ok := ob.UserOrders[userID]; ok {
		var remainingContracts []Contract
		for _, contract := range contracts {
			if contract.OrderType == orderType && contract.Price == price {
				switch orderType {
				case "sell":
					ob.Asks[price] -= contract.Quantity
					if ob.Asks[price] == 0 {
						delete(ob.Asks, price)
					}
				case "buy":
					ob.Bids[price] -= contract.Quantity
					if ob.Bids[price] == 0 {
						delete(ob.Bids, price)
					}
				}
			} else {
				remainingContracts = append(remainingContracts, contract)
			}
		}
		ob.UserOrders[userID] = remainingContracts
	}
}

// GetOrderBook returns the top 5 asks and bids
func (ob *OrderBook) GetOrderBook() map[string]interface{} {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	return map[string]interface{}{
		"asks":              ob.getTopAsks(),
		"bids":              ob.getTopBids(),
		"last_matched_price": ob.LastMatchedPrice,
	}
}

// getTopAsks returns the top 5 asks
func (ob *OrderBook) getTopAsks() []map[string]interface{} {
	var asks []map[string]interface{}
	var prices []float64
	for price := range ob.Asks {
		prices = append(prices, price)
	}
	sort.Float64s(prices)
	for i := 0; i < min(5, len(prices)); i++ {
		price := prices[i]
		asks = append(asks, map[string]interface{}{"price": price, "quantity": ob.Asks[price]})
	}
	return asks
}

// getTopBids returns the top 5 bids
func (ob *OrderBook) getTopBids() []map[string]interface{} {
	var bids []map[string]interface{}
	var prices []float64
	for price := range ob.Bids {
		prices = append(prices, price)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(prices)))
	for i := 0; i < min(5, len(prices)); i++ {
		price := prices[i]
		bids = append(bids, map[string]interface{}{"price": price, "quantity": ob.Bids[price]})
	}
	return bids
}

// getLowestAskPrice returns the lowest ask price
func (ob *OrderBook) getLowestAskPrice() float64 {
	var prices []float64
	for price := range ob.Asks {
		prices = append(prices, price)
	}
	sort.Float64s(prices)
	return prices[0]
}

// getHighestBidPrice returns the highest bid price
func (ob *OrderBook) getHighestBidPrice() float64 {
	var prices []float64
	for price := range ob.Bids {
		prices = append(prices, price)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(prices)))
	return prices[0]
}

// Utility function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var orderBook = NewOrderBook()

func main() {
	r := gin.Default()

	// Route for creating an order
	r.POST("/order", func(c *gin.Context) {
		var contract Contract
		if err := c.BindJSON(&contract); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		contract.Timestamp = time.Now().Unix()
		orderBook.AddContract(contract)
		c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
	})

	// Route for canceling an order
	r.DELETE("/order", func(c *gin.Context) {
		var data struct {
			UserID    string  `json:"user_id"`
			OrderType string  `json:"order_type"`
			Price     float64 `json:"price"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		orderBook.CancelContract(data.UserID, data.OrderType, data.Price)
		c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
	})

	// Route for retrieving the order book
	r.GET("/orderbook", func(c *gin.Context) {
		c.JSON(http.StatusOK, orderBook.GetOrderBook())
	})

	// Run the server on port 8080
	r.Run(":8080")
}
