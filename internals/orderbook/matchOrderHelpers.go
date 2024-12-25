package orderbook

import (
	"log"
	"time"
)

// AddLimitOrdersToOrderBook adds all the orders that exists in the limit order tracker into the main ask and buy heaps if there are any to be added
func (ob *OrderBook) AddLimitOrdersToOrderBook() error {

	lowestAsk, err := ob.Asks.Top()

	if err != nil {
		return err
	}
	if ob.LimitAsks.TopPrice() != -1 {
		for ob.LimitAsks.TopPrice() <= lowestAsk.GetPrice() {
			topLimitAsks, err := ob.LimitAsks.Pop()
			if err != nil {
				return err
			}
			topLimitAsks.SetOrderType("sell")
			topLimitAsks.SetTimestamp(time.Now().UnixMilli())
			ob.AddContract(topLimitAsks)
		}
	}
	highestBid, err := ob.Bids.Top()

	if err != nil {
		return err
	}

	if ob.LimitBids.TopPrice() != -1 {
		for ob.LimitBids.TopPrice() >= highestBid.GetPrice() {
			topLimitBids, err := ob.LimitBids.Pop()
			if err != nil {
				return err
			}
			topLimitBids.SetOrderType("buy")
			topLimitBids.SetTimestamp(time.Now().UnixMilli())
			ob.AddContract(topLimitBids)
		}
	}
	return nil
}

// MergeTopPrices takes and merges contracts in the top of the ask books and the bid books
func (ob *OrderBook) MergeTopPrices() {
	lowestAskPrice := ob.Asks.TopPrice()
	highestBidPrice := ob.Bids.TopPrice()

	if lowestAskPrice == -1 || highestBidPrice == -1 {
		return
	}

	for ob.Asks.TopPrice() <= ob.Bids.TopPrice() {

		if ob.Asks.TopPrice() == -1 || ob.Bids.TopPrice() == -1 {
			break
		}
		_lowestAsk, err1 := ob.Asks.Top()
		_highestBid, err2 := ob.Bids.Top()
		if err1 != nil || err2 != nil {
			break
		}

		if _lowestAsk.GetQuantity() == _highestBid.GetQuantity() {
			lowestAsk, _ := ob.Asks.Pop()
			highestBid, _ := ob.Bids.Pop()
			ob.LogHandler(lowestAsk, highestBid)
		} else if _lowestAsk.GetQuantity() < _highestBid.GetQuantity() {
			lowestAsk, _ := ob.Asks.Pop()
			highestBid, _ := ob.Bids.Top()
			highestBid.SetQuantity(highestBid.GetQuantity() - lowestAsk.GetQuantity())
			ob.LogHandler(lowestAsk, highestBid)
		} else {
			lowestAsk, _ := ob.Asks.Top()
			highestBid, _ := ob.Bids.Pop()
			lowestAsk.SetQuantity(lowestAsk.GetQuantity() - highestBid.GetQuantity())
			ob.LogHandler(lowestAsk, highestBid)
		}
		log.Printf("Matched %v with %v", _lowestAsk.GetPrice(), _highestBid.GetPrice())
	}
}
