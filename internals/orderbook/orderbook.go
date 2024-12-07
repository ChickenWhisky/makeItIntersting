package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/helpers"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/emirpasic/gods/queues/arrayqueue"
	"github.com/emirpasic/gods/queues/priorityqueue"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
)

// OrderBook stores order data and handles order processing.

type Level struct {
	Price         float32                    // Price at that Level
	Type          bool                       // Ask(0) or Bid(1) for setting up the comparator for the hashmap
	NoOfContracts int32                      // If 0 then simply delete struct from parent hashmap
	Orders        *priorityqueue.Queue       // Meant to keep order of contracts based on TimeStamps
	Contracts     map[string]models.Contract // Keep track of contracts in order to get instant data
	ToBeDeleted   arrayqueue.Queue           // For tracking what needs to be deleted
}

type OrderBook struct {
	AsksLevelByLevel             *priorityqueue.Queue             // To keep track of orders simply by price level and not on a order basis
	BidsLevelByLevel             *priorityqueue.Queue             // Same except for Bids
	AsksOrderByOrder             map[float32]*priorityqueue.Queue // To be able to extract orders on a contract by contract basis
	BidsOrderByOrder             map[float32]*priorityqueue.Queue // Same except for Bids
	LimitOrderAsksLevelByLevel   *priorityqueue.Queue             // To keep track of Limit Orders based ordered by Prices and further by time stamps so
	LimitOrderBidsLevelByLevel   *priorityqueue.Queue             //
	AsksOrderByOrderLevelByLevel map[float32]*priorityqueue.Queue // To be able to extract orders on a contract by contract basis
	BidsOrderByOrderLevelByLevel map[float32]*priorityqueue.Queue // Same except for Bids
	IncomingContracts            chan models.Contract             // Channel to stream incoming orders
	UserOrders                   map[string][]models.Contract     // A map to extract any existing order
	ToBeDeletedOrders            arrayqueue.Queue                 // A map to keep track of
	LastMatchedPrices            []float32
	mu                           sync.Mutex
}

// NewOrderBook creates a new empty order book.

// Will try to implement order book in a better manner

func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		Asks:                       priorityqueue.NewWith(ForAsks),
		Bids:                       priorityqueue.NewWith(ForBids),
		LimitOrderAsksLevelByLevel: priorityqueue.NewWith(ForLimitOrdersAsk),
		LimitOrderBidsLevelByLevel: priorityqueue.NewWith(ForLimitOrdersBid),
		IncomingContracts:          make(chan models.Contract),
		UserOrders:                 make(map[string][]models.Contract),
		LastMatchedPrices:          make([]float32, 0),
	}
	ob.StartProcessing()
	return ob
}

//  Starts up a thread that continuously checks the channel for new contracts to process

func (ob *OrderBook) StartProcessing() {
	go func() {
		for contract := range ob.IncomingContracts {
			switch contract.RequestType {
			case "add":
				ob.AddContract(contract)
			case "delete":
				ob.CancelContract(contract)
			}
		}
	}()
}

// AddContract adds a new contract (order) to the order book and attempts to match orders.

func (ob *OrderBook) AddContract(contract models.Contract) {

	//ob.mu.Lock()
	//defer ob.mu.Unlock()

	switch contract.OrderType {
	case "sell":
		{

			lengthFromEnv, err := strconv.Atoi(os.Getenv("CONTRACT_ID_LENGTH"))
			if err != nil {
				log.Println("Error converting CONTRACT_ID_LENGTH to int")
			}
			contract.Contract_ID = helpers.GenerateRandomString(lengthFromEnv)
			ob.Asks.Enqueue(contract)
		}
	case "buy":
		{

			lengthFromEnv, err := strconv.Atoi(os.Getenv("CONTRACT_ID_LENGTH"))
			if err != nil {
				log.Println("Error converting CONTRACT_ID_LENGTH to int")
			}
			contract.Contract_ID = helpers.GenerateRandomString(lengthFromEnv)
			ob.Bids.Enqueue(contract)
		}
	case "limit_buy":
		ob.LimitOrderBidsLevelByLevel.Enqueue(contract)
	case "limit_sell":
		ob.LimitOrderAsksLevelByLevel.Enqueue(contract)
	}

	// Store the user's orders
	ob.UserOrders[contract.UserID] = append(ob.UserOrders[contract.UserID], contract)

	// Attempt to match orders after adding a new one
	ob.matchOrders()
}

// CancelContract cancels a specific user's contract.
func (ob *OrderBook) CancelContract(contract models.Contract) {

	//ob.mu.Lock()
	//defer ob.mu.Unlock()

	userID := contract.UserID

	if contracts, ok := ob.UserOrders[userID]; ok {
		var remainingContracts []models.Contract
		for _, contract := range contracts {
			if contract.OrderType == orderType && contract.Price == price {
				switch orderType {
				case "sell":
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
func getLastNElements(slice []float32, n int) []float32 {
	if n > len(slice) {
		n = len(slice)
	}
	return slice[len(slice)-n:]
}

// matchOrders matches the highest bid with the lowest ask.
func (ob *OrderBook) matchOrders() {

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
			ob.removeLimitOrder(&topAsk.Price, topAsk, ob.LimitOrderAsksLevelByLevel)
		}
		if topBid.Quantity == 0 {
			ob.Bids.Dequeue()
			ob.removeLimitOrder(&topBid.Price, topBid, ob.LimitOrderBidsLevelByLevel)
		}
	}
}

// Helper to remove a specific order from limit order queue
func (ob *OrderBook) removeLimitOrder(price *float32, contract *models.Contract, limitOrders map[float32]*priorityqueue.Queue) {
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
func (ob *OrderBook) getTopAsks() []map[string]interface{} {
	var asks []map[string]interface{}
	var prices []float32
	for price := range ob.Asks {
		prices = append(prices, price)
	}
	sort.float32s(prices)
	for i := 0; i < min(5, len(prices)); i++ {
		price := prices[i]
		asks = append(asks, map[string]interface{}{"price": price, "quantity": ob.Asks[price]})
	}
	return asks
}

// getTopBids returns the top 5 bids with total quantity.
func (ob *OrderBook) getTopBids() []map[string]interface{} {
	var bids []map[string]interface{}
	var prices []float32
	for price := range ob.Bids {
		prices = append(prices, price)
	}
	sort.Sort(sort.Reverse(sort.float32Slice(prices)))
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
