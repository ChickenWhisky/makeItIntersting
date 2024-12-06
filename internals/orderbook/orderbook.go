package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/emirpasic/gods/queues/priorityqueue"
	"log"
	"sort"
)

// NewOrderBook creates a new empty order book.

// Will try to implement order book in a better manner

func NewOrderBook() *models.OrderBook {
	return &models.OrderBook{
		Asks:              priorityqueue.NewWith(ForAsks),
		Bids:              priorityqueue.NewWith(ForBids),
		LimitOrderAsks:    priorityqueue.NewWith(ForLimitOrdersAsk),
		LimitOrderBids:    priorityqueue.NewWith(ForLimitOrdersBid),
		IncomingContracts: make(chan models.Contract),
		UserOrders:        make(map[string][]models.Contract),
		LastMatchedPrices: make([]float64, 0),
	}
}

// AddContract adds a new contract (order) to the order book and attempts to match orders.

func (ob *models.OrderBook) AddContract(contract models.Contract) {

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
func (ob *models.OrderBook) CancelContract(userID, orderType string, price float64) {
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
func (ob *models.OrderBook) GetOrderBook() map[string]interface{} {
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
func (ob *models.OrderBook) matchOrders() {

	for ob.Asks.Size() > 0 && ob.Bids.Size() > 0 {
		// Peek top ask and bid
		topAskInterface, askOk := ob.Asks.Peek()
		topBidInterface, bidOk := ob.Bids.Peek()

		if !askOk || !bidOk {
			break
		}

		topAsk := topAskInterface.(*models.Contract)
		topBid := topBidInterface.(*models.Contract)

		// Determine matching quantity
		matchQuantity := min(topAsk.Quantity, topBid.Quantity)

		// Update order quantities
		topAsk.Quantity -= matchQuantity
		topBid.Quantity -= matchQuantity

		// Record matched price
		ob.LastMatchedPrices = append(ob.LastMatchedPrices, topAsk.Price)

		// Remove fully matched orders
		if topAsk.Quantity == 0 {
			ob.Asks.Dequeue()
			ob.removeLimitOrder(&topAsk.Price, topAsk, ob.LimitOrderAsks)
		}
		if topBid.Quantity == 0 {
			ob.Bids.Dequeue()
			ob.removeLimitOrder(&topBid.Price, topBid, ob.LimitOrderBids)
		}
	}
}

// Helper to remove a specific order from limit order queue
func (ob *models.OrderBook) removeLimitOrder(price *float64, contract *models.Contract, limitOrders map[float64]*priorityqueue.Queue) {
	// Get the queue for this price level
	priceQueue := limitOrders[*price]

	// Create a new queue to store remaining orders
	newQueue := priorityqueue.New()

	// Iterate through the existing queue
	for !priceQueue.Empty() {
		currentContract, _ := priceQueue.Dequeue()

		// Add back all contracts except the one to be removed
		if currentContract.(*models.Contract) != contract {
			newQueue.Enqueue(currentContract)
		}
	}

	// Replace the old queue with the new one
	if newQueue.Empty() {
		delete(limitOrders, *price)
	} else {
		limitOrders[*price] = newQueue
	}
}

// getLowestAskPrice returns the lowest ask price.

// getHighestBidPrice returns the highest bid price.

// getTopAsks returns the top 5 asks with total quantity.
func (ob *models.OrderBook) getTopAsks() []map[string]interface{} {
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
func (ob *models.OrderBook) getTopBids() []map[string]interface{} {
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
