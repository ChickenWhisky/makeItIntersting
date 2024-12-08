package orderbook

import (
	"errors"
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
	LevelID       string                      // LevelBook ID
	Price         float32                     // Price at that LevelBook
	Type          bool                        // Ask(true) or Bid(false) for setting up the comparator for the hashmap
	NoOfContracts int64                       // If 0 then simply delete struct from parent hashmap
	Orders        *priorityqueue.Queue        // Meant to keep order of contracts based on TimeStamps
	Contracts     map[string]*models.Contract // Keep track of contracts in order to get instant data
	ToBeDeleted   map[string]*models.Contract // For tracking what needs to be deleted
}

type OrderBook struct {
	AsksLevelByLevel      *priorityqueue.Queue        // To keep track of orders simply by price level and not on a per-order basis
	BidsLevelByLevel      *priorityqueue.Queue        // Same except for Bids
	LimitAsksLevelByLevel *priorityqueue.Queue        // To keep track of Limit Orders based ordered by Prices and further by time stamps so
	LimitBidsLevelByLevel *priorityqueue.Queue        // Same except for Limit Order Bids
	AsksOrderByOrder      map[float32]*LevelBook      // To be able to extract orders on a contract by contract basis
	BidsOrderByOrder      map[float32]*LevelBook      // Same except for Bids
	LimitAsksOrderByOrder map[float32]*LevelBook      // To be able to extract orders on a contract by contract basis
	LimitBidsOrderByOrder map[float32]*LevelBook      // Same except for Limit Order Bids
	IncomingContracts     chan models.Contract        // Channel to stream incoming orders
	Orders                map[string]*models.Contract // A map to extract any existing order
	//ToBeDeletedOrders     map[string]*models.Contract // A map to keep track of
	ToBeDeletedLevels map[string]*LevelBook // A map to keep track of Levels to be deleted
	LastMatchedPrices []float32
	Lock              sync.Mutex
}

// NewOrderBook creates a new empty order book.
func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		AsksLevelByLevel:      priorityqueue.NewWith(LevelByLevel),
		BidsLevelByLevel:      priorityqueue.NewWith(LevelByLevel),
		LimitAsksLevelByLevel: priorityqueue.NewWith(LevelByLevel),
		LimitBidsLevelByLevel: priorityqueue.NewWith(LevelByLevel),
		AsksOrderByOrder:      make(map[float32]*LevelBook),
		BidsOrderByOrder:      make(map[float32]*LevelBook),
		LimitAsksOrderByOrder: make(map[float32]*LevelBook),
		LimitBidsOrderByOrder: make(map[float32]*LevelBook),
		IncomingContracts:     make(chan models.Contract),
		//ToBeDeletedOrders:     make(map[string]*models.Contract),
		ToBeDeletedLevels: make(map[string]*LevelBook),
		Orders:            make(map[string]*models.Contract),
		LastMatchedPrices: make([]float32, 0),
		Lock:              sync.Mutex{},
	}
	ob.StartProcessing()
	return ob
}

// StartProcessing tarts up a thread that continuously checks the channel for new contracts to process
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

func (ob *OrderBook) PushContractIntoQueue(contract models.Contract) {
	ob.IncomingContracts <- contract
}

func (ob *OrderBook) AddContract(contract models.Contract) error {

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
	default:
		return errors.New("invalid order type")
	}

	if err == nil {
		return errors.New("could not create order")
	}
	// Store the user's orders
	ob.Orders[contract.ContractID] = &contract
	// Attempt to match orders after adding a new one
	ob.MatchOrders()
	return nil
}

// CancelContract cancels a specific user's contract.
func (ob *OrderBook) CancelContract(contract models.Contract) error {

	// Data required to cancel a given contract
	// Contract_ID
	// User_ID
	// From the Contract_ID alone we can get the info like price and from there we can
	// easily access the required LevelBook and then delete the data accordingly

	// First check if Order Exists
	orderInSystem, ok := ob.Orders[contract.ContractID]
	if ok && orderInSystem.UserID == contract.UserID {
		switch contract.OrderType {
		case "buy":
			ob.DeleteContractFromAsks(contract)
		case "sell":
			ob.DeleteContractFromBids(contract)
		case "limit_buy":
			ob.DeleteContractFromLimitAsks(contract)
		case "limit_sell":
			ob.DeleteContractFromLimitBids(contract)
		default:
			{
				log.Println("Invalid order type")
				return errors.New("invalid order type")
			}
		}

	} else if !ok {
		log.Println("Order doesnt exist in the system")
		return errors.New("order doesnt exist in the system")
	} else {
		log.Println("Order doesnt belong to the user :", contract.UserID)
		return errors.New("order doesnt belong to the user")
	}

	// Try matching remaining elements on deletion
	ob.MatchOrders()

	return nil

}

// ModifyContract cancels a specific user's contract and then adds a new contract based on the updated modifications.
func (ob *OrderBook) ModifyContract(contract models.Contract) {

	// Data required to cancel a given contract
	// Contract_ID
	// User_ID
}

// MatchOrders matches the highest bid with the lowest ask.
func (ob *OrderBook) MatchOrders() {

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
