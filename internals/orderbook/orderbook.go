package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/helpers"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/emirpasic/gods/queues/priorityqueue"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// OrderBook stores order data and handles order processing.

type LevelBook struct {
	Price         float32                     // Price at that LevelBook
	Type          bool                        // Ask(true) or Bid(false) for setting up the comparator for the hashmap
	NoOfContracts int64                       // If 0 then simply delete struct from parent hashmap
	Orders        *priorityqueue.Queue        // Meant to keep order of contracts based on TimeStamps
	Contracts     map[string]*models.Contract // Keep track of contracts in order to get instant data
	ToBeDeleted   chan models.Contract        // For tracking what needs to be deleted
}

type OrderBook struct {
	AsksLevelByLevel      *priorityqueue.Queue          // To keep track of orders simply by price level and not on a per-order basis
	BidsLevelByLevel      *priorityqueue.Queue          // Same except for Bids
	AsksOrderByOrder      map[float32]*LevelBook        // To be able to extract orders on a contract by contract basis
	BidsOrderByOrder      map[float32]*LevelBook        // Same except for Bids
	LimitAsksLevelByLevel *priorityqueue.Queue          // To keep track of Limit Orders based ordered by Prices and further by time stamps so
	LimitBidsLevelByLevel *priorityqueue.Queue          // Same except for Limit Order Bids
	LimitAsksOrderByOrder map[float32]*LevelBook        // To be able to extract orders on a contract by contract basis
	LimitBidsOrderByOrder map[float32]*LevelBook        // Same except for Limit Order Bids
	IncomingContracts     chan models.Contract          // Channel to stream incoming orders
	UserOrders            map[string][]*models.Contract // A map to extract any existing order
	ToBeDeletedOrders     chan models.Contract          // A map to keep track of
	LastMatchedPrices     []float32
	Lock                  sync.Mutex
}

// NewOrderBook creates a new empty order book.

// Will try to implement order book in a better manner

func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		AsksLevelByLevel:      priorityqueue.NewWith(LevelByLevel),
		BidsLevelByLevel:      priorityqueue.NewWith(LevelByLevel),
		AsksOrderByOrder:      make(map[float32]*LevelBook),
		BidsOrderByOrder:      make(map[float32]*LevelBook),
		LimitAsksLevelByLevel: priorityqueue.NewWith(LevelByLevel),
		LimitBidsLevelByLevel: priorityqueue.NewWith(LevelByLevel),
		LimitAsksOrderByOrder: make(map[float32]*LevelBook),
		LimitBidsOrderByOrder: make(map[float32]*LevelBook),
		IncomingContracts:     make(chan models.Contract),
		ToBeDeletedOrders:     make(chan models.Contract),
		UserOrders:            make(map[string][]*models.Contract),
		LastMatchedPrices:     make([]float32, 0),
		Lock:                  sync.Mutex{},
	}
	ob.StartProcessing()
	return ob
}

//  Starts up a thread that continuously checks the channel for new contracts to process

func (ob *OrderBook) StartProcessing() {
	go func() {
		for contract := range ob.IncomingContracts {
			contract.Timestamp = time.Now().UnixMilli()
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
	lengthFromEnv, err := strconv.Atoi(os.Getenv("CONTRACT_ID_LENGTH"))
	if err != nil {
		log.Println("Error converting CONTRACT_ID_LENGTH to int")
	}
	contract.ContractID = helpers.GenerateRandomString(lengthFromEnv)
	switch contract.OrderType {
	case "sell":
		ob.AddContractToAsks(contract)
	case "buy":
		ob.AddContractToAsks(contract)
	case "limit_buy":
		ob.AddContractToLimitAsks(contract)
	case "limit_sell":
		ob.AddContractToLimitBids(contract)
	}
	// Store the user's orders
	ob.UserOrders[contract.UserID] = append(ob.UserOrders[contract.UserID], &contract)
	// Attempt to match orders after adding a new one
	ob.matchOrders()
}

// CancelContract cancels a specific user's contract.
func (ob *OrderBook) CancelContract(contract models.Contract) {

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
			ob.removeLimitOrder(&topAsk.Price, topAsk, ob.LimitAsksLevelByLevel)
		}
		if topBid.Quantity == 0 {
			ob.Bids.Dequeue()
			ob.removeLimitOrder(&topBid.Price, topBid, ob.LimitBidsLevelByLevel)
		}
	}
}
