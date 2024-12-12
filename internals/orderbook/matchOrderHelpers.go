package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"time"
)

func (ob *OrderBook) FinalLevelDeletion() {
	for !ob.AsksLevelByLevel.Empty() {
		l, _ := ob.AsksLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.AsksLevelByLevel.Dequeue()
		} else {
			break
		}
	}
	for !ob.BidsLevelByLevel.Empty() {
		l, _ := ob.BidsLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.BidsLevelByLevel.Dequeue()
		} else {
			break
		}
	}
	for !ob.LimitAsksLevelByLevel.Empty() {
		l, _ := ob.LimitAsksLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.LimitAsksLevelByLevel.Dequeue()
		} else {
			break
		}
	}
	for !ob.LimitBidsLevelByLevel.Empty() {
		l, _ := ob.LimitBidsLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.LimitBidsLevelByLevel.Dequeue()
		} else {
			break
		}
	}

}

func (ob *OrderBook) FinalContractDeletion() {

	l, BidsBookEmpty := ob.BidsLevelByLevel.Peek()
	// Check if the Bids Heap is empty if it is then there is nothing to delete in it
	if !BidsBookEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			c, _ := Level.Orders.Peek()
			contract := c.(*models.Contract)
			_, isToBeDeleted := Level.ToBeDeleted[contract.ContractID]
			if isToBeDeleted {
				delete(Level.ToBeDeleted, contract.ContractID)
				Level.Orders.Dequeue()
				Level.NoOfContracts -= contract.Quantity
			} else {
				break
			}
		}
	}

	l, AsksBookEmpty := ob.AsksLevelByLevel.Peek()
	// Check if the Asks Heap is empty if it is then there is nothing to delete in it
	if !AsksBookEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			c, _ := Level.Orders.Peek()
			contract := c.(*models.Contract)
			_, isToBeDeleted := Level.ToBeDeleted[contract.ContractID]
			if isToBeDeleted {
				delete(Level.ToBeDeleted, contract.ContractID)
				Level.Orders.Dequeue()
				Level.NoOfContracts -= contract.Quantity
			} else {
				break
			}
		}
	}

	l, LimitBidsBookEmpty := ob.LimitBidsLevelByLevel.Peek()
	// Check if the LimitBids Heap is empty if it is then there is nothing to delete in it
	if !LimitBidsBookEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			c, _ := Level.Orders.Peek()
			contract := c.(*models.Contract)
			_, isToBeDeleted := Level.ToBeDeleted[contract.ContractID]
			if isToBeDeleted {
				delete(Level.ToBeDeleted, contract.ContractID)
				Level.Orders.Dequeue()
				Level.NoOfContracts -= contract.Quantity
			} else {
				break
			}
		}
	}

	l, LimitAsksBookEmpty := ob.LimitAsksLevelByLevel.Peek()
	// Check if the LimitAsks Heap is empty if it is then there is nothing to delete in it
	if !LimitAsksBookEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			c, _ := Level.Orders.Peek()
			contract := c.(*models.Contract)
			_, isToBeDeleted := Level.ToBeDeleted[contract.ContractID]
			if isToBeDeleted {
				delete(Level.ToBeDeleted, contract.ContractID)
				Level.Orders.Dequeue()
				Level.NoOfContracts -= contract.Quantity
			} else {
				break
			}
		}
	}
}

// AddLimitOrdersToOrderBook adds all the orders that exists in the limit order tracker into the main ask and buy heaps if there are any to be added
func (ob *OrderBook) AddLimitOrdersToOrderBook() {
	for ob.LimitAsksLevelByLevel.Empty() {
		lap, LimitAsksBookEmpty := ob.LimitAsksLevelByLevel.Peek()
		if !LimitAsksBookEmpty {
			LimitAskLevel := lap.(*LevelBook)
			al, _ := ob.AsksLevelByLevel.Peek()
			AskLevel := al.(*LevelBook)

			// Check if the Top price in the LimitBook is valid inorder to be added into the Ask Book
			// i.e. LimitAskLevel.Peek().Price<= AskLevel.Peek().Price
			if LimitAskLevel.Price <= AskLevel.Price {
				for price, Level := range ob.LimitAsksOrderByOrder {
					if price <= AskLevel.Price {
						// Add all the contracts from the LimitAsksOrderByOrder to the AsksOrderByOrder
						for Level.Orders.Empty() {
							c, _ := Level.Orders.Peek()
							contract := c.(*models.Contract)
							if Level.ToBeDeleted[contract.ContractID] != nil {
								contract.Timestamp = time.Now().UnixMilli()
								contract.OrderType = "sell"
								ob.AddContractToAsks(*contract)
							}
							Level.NoOfContracts -= contract.Quantity
							Level.Orders.Dequeue()
						}
						// Delete the LevelBook from the LimitAsksOrderByOrder
						delete(ob.LimitAsksOrderByOrder, price)
					}
				}
			} else {
				break
			}
		}
	}
}

// MergeTopPrices takes and merges contracts in the top of the ask books and the bid books
func (ob *OrderBook) MergeTopPrices() {

	// Check if the top most bids will match or not
	lal, doAsksExist := ob.AsksLevelByLevel.Peek()
	hbl, doBidsExist := ob.BidsLevelByLevel.Peek()
	lowestAskLevel := lal.(*LevelBook)
	highestBidLevel := hbl.(*LevelBook)

	if doAsksExist && doBidsExist {
		if lowestAskLevel.Price <= highestBidLevel.Price {
			for !lowestAskLevel.Orders.Empty() && !highestBidLevel.Orders.Empty() {
				lowestAskContract, _ := lowestAskLevel.Orders.Peek()
				highestBidContract, _ := highestBidLevel.Orders.Peek()

				if lowestAskContract.(*models.Contract).Quantity == highestBidContract.(*models.Contract).Quantity {
					lowestAskLevel.Orders.Dequeue()
					highestBidLevel.Orders.Dequeue()
					lowestAskLevel.NoOfContracts -= lowestAskContract.(*models.Contract).Quantity
					highestBidLevel.NoOfContracts -= highestBidContract.(*models.Contract).Quantity
					logHandler(lowestAskContract.(*models.Contract), highestBidContract.(*models.Contract))
				} else if lowestAskContract.(*models.Contract).Quantity < highestBidContract.(*models.Contract).Quantity {
					lowestAskLevel.Orders.Dequeue()
					highestBidContract.(*models.Contract).Quantity -= lowestAskContract.(*models.Contract).Quantity
					lowestAskLevel.NoOfContracts -= lowestAskContract.(*models.Contract).Quantity
					logHandler(lowestAskContract.(*models.Contract), highestBidContract.(*models.Contract))
				} else {
					highestBidLevel.Orders.Dequeue()
					lowestAskContract.(*models.Contract).Quantity -= highestBidContract.(*models.Contract).Quantity
					highestBidLevel.NoOfContracts -= highestBidContract.(*models.Contract).Quantity
					logHandler(lowestAskContract.(*models.Contract), highestBidContract.(*models.Contract))
				}
			}
		}
	} else {
		return
	}
}
