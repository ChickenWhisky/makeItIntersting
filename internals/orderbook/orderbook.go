package orderbook

import (
	"log"
	"sort"
	"sync"
	// pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	// "github.com/emirpasic/gods/utils"

)

// OrderBook stores order data and handles order processing.
type OrderBook struct {
	Asks              map[float64]int
	Bids              map[float64]int
	UserOrders        map[string][]models.Contract
	LastMatchedPrices []float64
	mu                sync.Mutex
}

// NewOrderBook creates a new empty order book.
func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks:              make(map[float64]int),
		Bids:              make(map[float64]int),
		UserOrders:        make(map[string][]models.Contract),
		LastMatchedPrices: make([]float64, 0),
	}
}

// AddContract adds a new contract (order) to the order book and attempts to match orders.
func (ob *OrderBook) AddContract(contract models.Contract) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	switch contract.OrderType {
	case "sell":
		ob.Asks[contract.Price] += contract.Quantity
		if ob.Asks[contract.Price] < 0 {
			log.Printf("ERROR: Negative quantity for price %f in Asks: %d", contract.Price, ob.Asks[contract.Price])
		}
	case "buy":
		ob.Bids[contract.Price] += contract.Quantity
		if ob.Bids[contract.Price] < 0 {
			log.Printf("ERROR: Negative quantity for price %f in Bids: %d", contract.Price, ob.Bids[contract.Price])
		}
	}

	// Store the user's orders
	ob.UserOrders[contract.UserID] = append(ob.UserOrders[contract.UserID], contract)

	// Attempt to match orders after adding a new one
	ob.matchOrders()
}

// CancelContract cancels a specific user's contract.
func (ob *OrderBook) CancelContract(userID, orderType string, price float64) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if contracts, ok := ob.UserOrders[userID]; ok {
		var remainingContracts []models.Contract
		for _, contract := range contracts {
			if contract.OrderType == orderType && contract.Price == price {
				switch orderType {
				case "sell":
					// if ob.Asks[price] >= contract.Quantity{
					// 	ob.Asks[price] -= contract.Quantity
					// }
					ob.Asks[price] -= contract.Quantity
					if ob.Asks[price] == 0 {
						delete(ob.Asks, price)
					}
					// else if ob.Asks[price] < 0{
					// 	fmt.Println("Smn wrong in line 63 orderbook.go")
					// }
				case "buy":
					ob.Bids[price] -= contract.Quantity
					// if ob.Bids[price] >= contract.Quantity{
					// 	ob.Bids[price] -= contract.Quantity
					// }
					if ob.Bids[price] == 0 {
						delete(ob.Bids, price)
					}
					// else if ob.Bids[price] < 0{
					// 	fmt.Println("Smn wrong in line 71 orderbook.go")
					// }
				}
			} else {
				remainingContracts = append(remainingContracts, contract)
			}
		}
		ob.UserOrders[userID] = remainingContracts
	}
}

// GetOrderBook returns the current state of the order book.
func (ob *OrderBook) GetOrderBook() map[string]interface{} {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	return map[string]interface{}{
		"asks":                   ob.getTopAsks(),
		"bids":                   ob.getTopBids(),
		"last_matched_price":     ob.LastMatchedPrices[len(ob.LastMatchedPrices)-1],
		"last_50_matched_prices": getLastNElements(ob.LastMatchedPrices, 100),
	}
}

// Helper methods

// getLastNElements returns the last n elements of a slice.
func getLastNElements(slice []float64, n int) []float64 {
	if n > len(slice) {
		n = len(slice)
	}
	return slice[len(slice)-n:]
}

// matchOrders matches the highest bid with the lowest ask.
func (ob *OrderBook) matchOrders() {
	for len(ob.Asks) > 0 && len(ob.Bids) > 0 {
		lowestAskPrice := ob.getLowestAskPrice()
		highestBidPrice := ob.getHighestBidPrice()

		if highestBidPrice < lowestAskPrice {
			break
		}

		matchedQuantity := min(ob.Asks[lowestAskPrice], ob.Bids[highestBidPrice])

		ob.LastMatchedPrices = append(ob.LastMatchedPrices, lowestAskPrice)
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

// getLowestAskPrice returns the lowest ask price.
func (ob *OrderBook) getLowestAskPrice() float64 {
	var prices []float64
	for price := range ob.Asks {
		prices = append(prices, price)
	}
	sort.Float64s(prices)
	return prices[0]
}

// getHighestBidPrice returns the highest bid price.
func (ob *OrderBook) getHighestBidPrice() float64 {
	var prices []float64
	for price := range ob.Bids {
		prices = append(prices, price)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(prices)))
	return prices[0]
}

// getTopAsks returns the top 5 asks with total quantity.
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

// getTopBids returns the top 5 bids with total quantity.
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

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
