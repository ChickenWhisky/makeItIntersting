package orderbook

import (
	"errors"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"strconv"
	"sync"
	"time"
)

// OrderBook stores order data and handles order processing.

type OrderBook struct {
	Asks              *models.OrderQueue          // OrderQueue for Asks
	Bids              *models.OrderQueue          // OrderQueue for Bids
	LimitAsks         *models.OrderQueue          // OrderQueue for LimitAsks
	LimitBids         *models.OrderQueue          // OrderQueue for LimitBids
	IncomingContracts chan models.Contract        // Channel to stream incoming orders
	CompletedTrades   chan models.Trade           // Channel to stream incoming orders
	Orders            map[string]*models.Contract // A map to extract any existing order
	LastMatchedPrices []models.Trade              // Struct to keep track of last matched prices
	Lock              sync.Mutex                  // Mutex (Maybe removed if not required)
	ContractNo        int                         // reference number to create ContractIDs
}

// NewOrderBook creates a new empty order book.
func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		Asks:              models.NewOrderQueue(),
		Bids:              models.NewOrderQueue(),
		LimitAsks:         models.NewOrderQueue(),
		LimitBids:         models.NewOrderQueue(),
		IncomingContracts: make(chan models.Contract),
		CompletedTrades:   make(chan models.Trade),
		Orders:            make(map[string]*models.Contract),
		LastMatchedPrices: make([]models.Trade, 0),
		Lock:              sync.Mutex{},
		ContractNo:        0,
	}
	ob.StartProcessingOrders()
	//ob.StartProcessingTrades()
	return ob
}

// StartProcessingTrades tarts up a thread that continuously checks the channel for new contracts to process
func (ob *OrderBook) StartProcessingTrades() {
	go func() {
		for contract := range ob.IncomingContracts {
			contract.SetTimestamp(time.Now().UnixMilli())
			log.Printf("Processing contract: %+v", contract)
			switch contract.RequestType {
			case "add":
				if err := ob.AddContract(&contract); err != nil {
					log.Printf("Error adding contract: %v", err)
				}
			case "delete":
				if err := ob.CancelContract(&contract); err != nil {
					log.Printf("Error deleting contract: %v", err)
				}
			default:
				log.Printf("Unknown request type: %s", contract.RequestType)
			}
		}
	}()
}

// StartProcessingOrders tarts up a thread that continuously checks the channel for new contracts to process
func (ob *OrderBook) StartProcessingOrders() {
	go func() {
		for contract := range ob.IncomingContracts {
			contract.SetTimestamp(time.Now().UnixMilli())
			log.Printf("Processing contract: %+v", contract)
			switch contract.RequestType {
			case "add":
				if err := ob.AddContract(&contract); err != nil {
					log.Printf("Error adding contract: %v", err)
				}
			case "delete":
				if err := ob.CancelContract(&contract); err != nil {
					log.Printf("Error deleting contract: %v", err)
				}
			case "modify":
				if err := ob.ModifyContract(&contract); err != nil {
					log.Printf("Error modifying contract: %v", err)
				}
			default:
				log.Printf("Unknown request type: %s", contract.RequestType)
			}
		}
	}()
}

func (ob *OrderBook) PushContractIntoQueue(contract models.Contract) {
	ob.IncomingContracts <- contract
}

// AddContract adds a new contract (order) to the order book and attempts to match orders.
func (ob *OrderBook) AddContract(contract *models.Contract) error {

	//contract.ContractID = helpers.GenerateRandomString(helpers.ConvertStringToInt(os.Getenv("CONTRACT_ID_LENGTH")))

	newContractID := strconv.Itoa(ob.ContractNo)
	contract.ContractID = newContractID
	//log.Infof("New Contract ID: %s ", contract.GetContractID())
	ob.ContractNo++
	switch contract.OrderType {
	case "sell":

		err := ob.Asks.Push(contract)
		if err != nil {
			log.Printf("Error pushing contract into Asks: %v", err)
		}
		log.Printf("Pushed contract into Asks: %+v", contract)
	case "buy":
		err := ob.Bids.Push(contract)
		if err != nil {
			log.Printf("Error pushing contract into Bids: %v", err)
		}
		log.Printf("Pushed contract into Bids: %+v", contract)
	case "limit_sell":
		err := ob.LimitAsks.Push(contract)
		if err != nil {
			log.Printf("Error pushing contract into LimitAsks: %v", err)
		}
		log.Printf("Pushed contract into LimitAsks: %+v", contract)
	case "limit_buy":
		err := ob.LimitBids.Push(contract)
		if err != nil {
			log.Printf("Error pushing contract into LimitBids: %v", err)
		}
		log.Printf("Pushed contract into LimitBids: %+v", contract)
	default:
		return errors.New("invalid order type")
	}

	// Store the user's orders
	ob.Orders[contract.ContractID] = contract
	// Attempt to match orders after adding a new one
	ob.MatchOrders()
	return nil
}

// CancelContract cancels a specific user's contract.
func (ob *OrderBook) CancelContract(contract *models.Contract) error {

	// Data required to cancel a given contract
	// Contract_ID
	// User_ID
	// From the Contract_ID alone we can get the info like price and from there we can
	// easily access the required LevelBook and then delete the data accordingly

	// First check if Order Exists
	orderInSystem, ok := ob.Orders[contract.ContractID]
	if ok && orderInSystem.UserID == contract.GetUserID() {
		log.Infof("Order found in the system : %v", contract.GetContractID())
		switch contract.OrderType {
		case "buy":
			ob.Asks.Delete(contract.GetContractID())
			_, err := ob.Asks.Find(contract.GetContractID())
			if err == nil {
				log.Printf("CONTRACT DELETION ERROR : %v", err)
			}
			log.Printf("CONTRACT DELETION SUCCESS")
		case "sell":
			ob.Bids.Delete(contract.GetContractID())
			_, err := ob.Bids.Find(contract.GetContractID())
			if err == nil {
				log.Printf("CONTRACT DELETION ERROR : %v", err)
			}
			log.Printf("CONTRACT DELETION SUCCESS")
		case "limit_buy":
			ob.LimitBids.Delete(contract.GetContractID())
			_, err := ob.LimitBids.Find(contract.GetContractID())
			if err == nil {
				log.Printf("CONTRACT DELETION ERROR : %v", err)
			}
			log.Printf("CONTRACT DELETION SUCCESS")
		case "limit_sell":
			ob.LimitAsks.Delete(contract.GetContractID())
			_, err := ob.LimitAsks.Find(contract.GetContractID())
			if err == nil {
				log.Printf("CONTRACT DELETION ERROR : %v", err)
			}
			log.Printf("CONTRACT DELETION SUCCESS")
		default:
			{
				return errors.New("invalid order type")
			}
		}
	} else if !ok {
		log.Printf("Order doesnt exist in the system")
		return errors.New("order doesnt exist in the system")
	} else {
		log.Printf("Order doesnt belong to the user :", contract.UserID)
		return errors.New("order doesnt belong to the user")
	}

	// Try matching remaining elements on deletion
	delete(ob.Orders, contract.GetContractID())
	ob.MatchOrders()

	return nil

}

// ModifyContract cancels a specific user's contract and then adds a new contract based on the updated modifications.
func (ob *OrderBook) ModifyContract(contract *models.Contract) error {

	switch contract.OrderType {
	case "buy":
		c, err := ob.Bids.Find(contract.ContractID)
		if err != nil {
			return err
		}
		c.SetPrice(contract.GetPrice())
		c.SetQuantity(contract.GetQuantity())
	case "sell":
		c, err := ob.Asks.Find(contract.ContractID)
		if err != nil {
			return err
		}
		c.SetPrice(contract.GetPrice())
		c.SetQuantity(contract.GetQuantity())
	case "limit_buy":
		c, err := ob.LimitBids.Find(contract.ContractID)
		if err != nil {
			return err
		}
		c.SetPrice(contract.GetPrice())
		c.SetQuantity(contract.GetQuantity())
	case "limit_sell":
		c, err := ob.LimitAsks.Find(contract.ContractID)
		if err != nil {
			return err
		}
		c.SetPrice(contract.GetPrice())
		c.SetQuantity(contract.GetQuantity())
	default:
		return errors.New("invalid order type")
	}
	ob.MatchOrders()
	return nil
}

// MatchOrders matches the highest bid with the lowest ask.
func (ob *OrderBook) MatchOrders() {

	err1 := ob.AddLimitOrdersToOrderBook()
	if err1 != nil {
		log.Printf("Error adding limit orders to order book: %v", err1)
	}
	err2 := ob.MergeTopPrices()
	if err2 != nil {
		log.Printf("Error merging top prices: %v", err2)
	}
}

// GetLastTrades gets the last n traded prices
func (ob *OrderBook) GetLastTrades(n int) []models.Trade {
	if n > len(ob.LastMatchedPrices) {
		n = len(ob.LastMatchedPrices)
	}
	return ob.LastMatchedPrices[len(ob.LastMatchedPrices)-n:]
}
