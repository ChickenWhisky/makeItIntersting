package orderbook

import (
	"errors"
	"github.com/charmbracelet/log"
	"time"
)

// AddLimitOrdersToOrderBook adds all the orders that exists in the limit order tracker into the main ask and buy heaps if there are any to be added
func (ob *OrderBook) AddLimitOrdersToOrderBook() error {

	//lowestAsk, err1 := ob.Asks.Top()
	log.Print("Starting Function: AddLimitOrdersToOrderBook")
	lowestAsk, _ := ob.Asks.Top()
	if lowestAsk != nil {
		if ob.LimitAsks.TopPrice() != -1 {
			for ob.LimitAsks.TopPrice() <= lowestAsk.GetPrice() {
				topLimitAsks, err := ob.LimitAsks.Pop()
				if err != nil {

					return errors.New("Error popping limit asks")
				}
				topLimitAsks.SetOrderType("sell")
				topLimitAsks.SetTimestamp(time.Now().UnixMilli())
				ob.AddOrder(topLimitAsks)
			}
		}
	}
	highestBid, _ := ob.Bids.Top()
	//highestBid, err2 := ob.Bids.Top()
	if highestBid != nil {
		if ob.LimitBids.TopPrice() != -1 {
			for ob.LimitBids.TopPrice() >= highestBid.GetPrice() {
				topLimitBids, err := ob.LimitBids.Pop()
				if err != nil {
					return errors.New("Error popping limit bids")
				}
				topLimitBids.SetOrderType("buy")
				topLimitBids.SetTimestamp(time.Now().UnixMilli())
				ob.AddOrder(topLimitBids)
			}
		}
	}
	log.Print("Ending Function: AddLimitOrdersToOrderBook")
	return nil
}

// MergeTopPrices takes and merges Orders in the top of the ask books and the bid books
func (ob *OrderBook) MergeTopPrices() error {
	log.Info("Starting Function: MergeTopPrices")
	lowestAskPrice := ob.Asks.TopPrice()

	highestBidPrice := ob.Bids.TopPrice()

	if lowestAskPrice == -1 || highestBidPrice == -1 {
		log.Infof("Ask Price: %v, Bid Price: %v", lowestAskPrice, highestBidPrice)
		return errors.New("No Orders to merge")
	}

	log.Infof("Ask Price: %v, Bid Price: %v", lowestAskPrice, highestBidPrice)
	for ob.Asks.TopPrice() <= ob.Bids.TopPrice() {
		log.Infof("Ask Price: %v, Bid Price: %v", ob.Asks.TopPrice(), ob.Bids.TopPrice())
		if ob.Asks.TopPrice() == -1 || ob.Bids.TopPrice() == -1 {
			break
		}
		_lowestAsk, err1 := ob.Asks.Top()
		_highestBid, err2 := ob.Bids.Top()
		if err1 != nil || err2 != nil {
			return err1
		}
		var NoOfOrders int64 = 0
		if _lowestAsk.GetQuantity() == _highestBid.GetQuantity() {
			lowestAsk, _ := ob.Asks.Pop()
			highestBid, _ := ob.Bids.Pop()
			NoOfOrders = _lowestAsk.GetQuantity()
			ob.LogHandler(lowestAsk, highestBid)
		} else if _lowestAsk.GetQuantity() < _highestBid.GetQuantity() {
			lowestAsk, _ := ob.Asks.Pop()
			highestBid, _ := ob.Bids.Top()
			NoOfOrders = lowestAsk.GetQuantity()
			ob.LogHandler(lowestAsk, highestBid)
			highestBid.SetQuantity(highestBid.GetQuantity() - lowestAsk.GetQuantity())
			delete(ob.Orders, lowestAsk.GetOrderID())
		} else {
			lowestAsk, _ := ob.Asks.Top()
			highestBid, _ := ob.Bids.Pop()
			NoOfOrders = highestBid.GetQuantity()
			ob.LogHandler(lowestAsk, highestBid)
			lowestAsk.SetQuantity(lowestAsk.GetQuantity() - highestBid.GetQuantity())
			delete(ob.Orders, highestBid.GetOrderID())
		}
		log.Printf("\nAsk_Price : %v\n Bid_Price : %v\n Orders : %v", _lowestAsk.GetPrice(), _highestBid.GetPrice(), NoOfOrders)
	}
	log.Info("Ending Function: MergeTopPrices")
	return nil
}
