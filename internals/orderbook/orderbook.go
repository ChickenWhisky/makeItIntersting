package orderbook

import (
	"errors"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/charmbracelet/log"
	"strconv"
	"sync"
	"time"
)

// OrderBook stores order data and handles order processing.

type OrderBook struct {
	Asks              *models.OrderQueue       // OrderQueue for Asks
	Bids              *models.OrderQueue       // OrderQueue for Bids
	LimitAsks         *models.OrderQueue       // OrderQueue for LimitAsks
	LimitBids         *models.OrderQueue       // OrderQueue for LimitBids
	IncomingOrders    chan models.Order        // Channel to stream incoming orders
	CompletedTrades   chan models.Trade        // Channel to stream incoming orders
	Orders            map[string]*models.Order // A map to extract any existing order
	LastMatchedPrices []models.Trade           // Struct to keep track of last matched prices
	Lock              sync.Mutex               // Mutex (Maybe removed if not required)
	OrderNo           int                      // reference number to create OrderIDs
	TradeNo           int                      // reference number to create TradeIDs
}

// NewOrderBook creates a new empty order book.
func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		Asks:              models.NewOrderQueue(),
		Bids:              models.NewOrderQueue(),
		LimitAsks:         models.NewOrderQueue(),
		LimitBids:         models.NewOrderQueue(),
		IncomingOrders:    make(chan models.Order),
		CompletedTrades:   make(chan models.Trade),
		Orders:            make(map[string]*models.Order),
		LastMatchedPrices: make([]models.Trade, 0),
		Lock:              sync.Mutex{},
		OrderNo:           0,
		TradeNo:           0,
	}
	ob.StartProcessingOrders()
	//ob.StartProcessingTrades()
	return ob
}

// StartProcessingTrades tarts up a thread that continuously checks the channel for new Orders to process
func (ob *OrderBook) StartProcessingTrades() {
	go func() {
		for Order := range ob.IncomingOrders {
			Order.SetTimestamp(time.Now().UnixMilli())
			log.Printf("Processing Order: %+v", Order)
			switch Order.RequestType {
			case "add":
				if err := ob.AddOrder(&Order); err != nil {
					log.Printf("Error adding Order: %v", err)
				}
			case "delete":
				if err := ob.CancelOrder(&Order); err != nil {
					log.Printf("Error deleting Order: %v", err)
				}
			default:
				log.Printf("Unknown request type: %s", Order.RequestType)
			}
		}
	}()
}

// StartProcessingOrders starts up a thread that continuously checks the channel for new Orders to process
func (ob *OrderBook) StartProcessingOrders() {
	go func() {
		for Order := range ob.IncomingOrders {
			Order.SetTimestamp(time.Now().UnixMilli())
			log.Printf("Processing Order: %+v", Order)
			switch Order.RequestType {
			case "add":
				if err := ob.AddOrder(&Order); err != nil {
					log.Printf("Error adding Order: %v", err)
				}
			case "delete":
				if err := ob.CancelOrder(&Order); err != nil {
					log.Printf("Error deleting Order: %v", err)
				}
			case "modify":
				if err := ob.ModifyOrder(&Order); err != nil {
					log.Printf("Error modifying Order: %v", err)
				}
			default:
				log.Printf("Unknown request type: %s", Order.RequestType)
			}
		}
	}()
}

func (ob *OrderBook) PushOrderIntoQueue(Order models.Order) {
	ob.IncomingOrders <- Order
}

// AddOrder adds a new Order (order) to the order book and attempts to match orders.
func (ob *OrderBook) AddOrder(Order *models.Order) error {

	//Order.OrderID = helpers.GenerateRandomString(helpers.ConvertStringToInt(os.Getenv("Order_ID_LENGTH")))

	newOrderID := strconv.Itoa(ob.OrderNo)
	Order.OrderID = newOrderID
	ob.OrderNo++
	switch Order.OrderType {
	case "sell":

		err := ob.Asks.Push(Order)
		if err != nil {
			log.Printf("Error pushing Order into Asks: %v", err)
		}
		log.Printf("Pushed Order into Asks: %+v", Order)
	case "buy":
		err := ob.Bids.Push(Order)
		if err != nil {
			log.Printf("Error pushing Order into Bids: %v", err)
		}
		log.Printf("Pushed Order into Bids: %+v", Order)
	case "limit_sell":
		err := ob.LimitAsks.Push(Order)
		if err != nil {
			log.Printf("Error pushing Order into LimitAsks: %v", err)
		}
		log.Printf("Pushed Order into LimitAsks: %+v", Order)
	case "limit_buy":
		err := ob.LimitBids.Push(Order)
		if err != nil {
			log.Printf("Error pushing Order into LimitBids: %v", err)
		}
		log.Printf("Pushed Order into LimitBids: %+v", Order)
	default:
		return errors.New("invalid order type")
	}

	// Store the user's orders
	ob.Orders[Order.OrderID] = Order
	// Attempt to match orders after adding a new one
	ob.MatchOrders()
	return nil
}

// CancelOrder cancels a specific user's Order.
func (ob *OrderBook) CancelOrder(Order *models.Order) error {

	// Data required to cancel a given Order
	// Order_ID
	// User_ID
	// From the Order_ID alone we can get the info like price and from there we can
	// easily access the required LevelBook and then delete the data accordingly

	// First check if Order Exists
	orderInSystem, ok := ob.Orders[Order.OrderID]
	if ok && orderInSystem.UserID == Order.GetUserID() {
		log.Infof("Order found in the system : %v", Order.GetOrderID())
		switch Order.OrderType {
		case "buy":
			ob.Asks.Delete(Order.GetOrderID())
			_, err := ob.Asks.Find(Order.GetOrderID())
			if err == nil {
				log.Printf("Order DELETION ERROR : %v", err)
			}
			log.Printf("Order DELETION SUCCESS")
		case "sell":
			ob.Bids.Delete(Order.GetOrderID())
			_, err := ob.Bids.Find(Order.GetOrderID())
			if err == nil {
				log.Printf("Order DELETION ERROR : %v", err)
			}
			log.Printf("Order DELETION SUCCESS")
		case "limit_buy":
			ob.LimitBids.Delete(Order.GetOrderID())
			_, err := ob.LimitBids.Find(Order.GetOrderID())
			if err == nil {
				log.Printf("Order DELETION ERROR : %v", err)
			}
			log.Printf("Order DELETION SUCCESS")
		case "limit_sell":
			ob.LimitAsks.Delete(Order.GetOrderID())
			_, err := ob.LimitAsks.Find(Order.GetOrderID())
			if err == nil {
				log.Printf("Order DELETION ERROR : %v", err)
			}
			log.Printf("Order DELETION SUCCESS")
		default:
			{
				return errors.New("invalid order type")
			}
		}
	} else if !ok {
		log.Printf("Order doesnt exist in the system")
		return errors.New("order doesnt exist in the system")
	} else {
		log.Printf("Order doesnt belong to the user : %v", Order.UserID)
		return errors.New("order doesnt belong to the user")
	}

	// Try matching remaining elements on deletion
	delete(ob.Orders, Order.GetOrderID())
	ob.MatchOrders()

	return nil

}

// ModifyOrder cancels a specific user's Order and then adds a new Order based on the updated modifications.
func (ob *OrderBook) ModifyOrder(Order *models.Order) error {

	switch Order.OrderType {
	case "buy":
		c, err := ob.Bids.Find(Order.OrderID)
		if err != nil {
			return err
		}
		c.SetPrice(Order.GetPrice())
		c.SetQuantity(Order.GetQuantity())
	case "sell":
		c, err := ob.Asks.Find(Order.OrderID)
		if err != nil {
			return err
		}
		c.SetPrice(Order.GetPrice())
		c.SetQuantity(Order.GetQuantity())
	case "limit_buy":
		c, err := ob.LimitBids.Find(Order.OrderID)
		if err != nil {
			return err
		}
		c.SetPrice(Order.GetPrice())
		c.SetQuantity(Order.GetQuantity())
	case "limit_sell":
		c, err := ob.LimitAsks.Find(Order.OrderID)
		if err != nil {
			return err
		}
		c.SetPrice(Order.GetPrice())
		c.SetQuantity(Order.GetQuantity())
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
	if n == -1 {
		return ob.LastMatchedPrices
	}
	if n > len(ob.LastMatchedPrices) {
		n = len(ob.LastMatchedPrices)
	}
	return ob.LastMatchedPrices[len(ob.LastMatchedPrices)-n:]
}
