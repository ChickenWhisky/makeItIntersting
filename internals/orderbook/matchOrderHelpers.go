package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"log"
	"time"
)

func (ob *OrderBook) FinalLevelDeletion() {
	for !ob.AsksLevelByLevel.Empty() {
		log.Printf("Loop in FinalLevelDeletion Asks")
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
		log.Printf("Loop in FinalLevelDeletion Bids")
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
		log.Printf("Loop in FinalLevelDeletion Limit Asks")
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
		log.Printf("Loop in FinalLevelDeletion Limit Bids")
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

	l, BidsBookNotEmpty := ob.BidsLevelByLevel.Peek()
	// Check if the Bids Heap is empty if it is then there is nothing to delete in it
	if BidsBookNotEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			log.Printf("Loop in FinalContractDeletion Bids")
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

	l, AsksBookNotEmpty := ob.AsksLevelByLevel.Peek()
	// Check if the Asks Heap is empty if it is then there is nothing to delete in it
	if AsksBookNotEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			log.Printf("Loop in FinalContractDeletion Asks")
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

	l, LimitBidsBookNotEmpty := ob.LimitBidsLevelByLevel.Peek()
	// Check if the LimitBids Heap is empty if it is then there is nothing to delete in it
	if LimitBidsBookNotEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			log.Printf("Loop in FinalContractDeletion LimitBids")
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

	l, LimitAsksBookNotEmpty := ob.LimitAsksLevelByLevel.Peek()
	// Check if the LimitAsks Heap is empty if it is then there is nothing to delete in it
	if LimitAsksBookNotEmpty {
		Level := l.(*LevelBook)
		for !Level.Orders.Empty() {
			log.Printf("Loop in FinalContractDeletion LimitAsks")
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
	for !ob.LimitAsksLevelByLevel.Empty() {
		log.Printf("Loop in AddLimitOrdersToOrderBook LimitAsks")
		lap, LimitAsksBookNotEmpty := ob.LimitAsksLevelByLevel.Peek()
		if LimitAsksBookNotEmpty {
			LimitAskLevel := lap.(*LevelBook)
			al, AskBookNotEmpty := ob.AsksLevelByLevel.Peek()
			if !AskBookNotEmpty {
				break
			}
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
								ob.AddContract(*contract)
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
	for !ob.LimitBidsLevelByLevel.Empty() {
		log.Printf("Loop in AddLimitOrdersToOrderBook LimitBids")
		lbp, LimitBidsBookNotEmpty := ob.LimitBidsLevelByLevel.Peek()
		if LimitBidsBookNotEmpty {
			LimitBidsLevel := lbp.(*LevelBook)
			bl, BidBookNotEmpty := ob.BidsLevelByLevel.Peek()
			if !BidBookNotEmpty {
				break
			}
			BidsLevel := bl.(*LevelBook)

			// Check if the Top price in the LimitBook is valid inorder to be added into the Ask Book
			// i.e. LimitAskLevel.Peek().Price<= AskLevel.Peek().Price
			if LimitBidsLevel.Price >= BidsLevel.Price {
				for price, Level := range ob.LimitBidsOrderByOrder {
					if price >= BidsLevel.Price {
						// Add all the contracts from the LimitAsksOrderByOrder to the AsksOrderByOrder
						for Level.Orders.Empty() {
							c, _ := Level.Orders.Peek()
							contract := c.(*models.Contract)
							if Level.ToBeDeleted[contract.ContractID] != nil {
								contract.Timestamp = time.Now().UnixMilli()
								contract.OrderType = "sell"
								ob.AddContract(*contract)
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
	lal, doAsksExist := ob.AsksLevelByLevel.Peek()
	hbl, doBidsExist := ob.BidsLevelByLevel.Peek()

	if !doAsksExist || !doBidsExist {
		return
	}

	lowestAskLevel := lal.(*LevelBook)
	highestBidLevel := hbl.(*LevelBook)

	for !lowestAskLevel.Orders.Empty() && !highestBidLevel.Orders.Empty() {
		lowestAskContract, _ := lowestAskLevel.Orders.Peek()
		highestBidContract, _ := highestBidLevel.Orders.Peek()
		lac := lowestAskContract.(*models.Contract)
		hbc := highestBidContract.(*models.Contract)

		if lowestAskLevel.ToBeDeleted[lac.ContractID] == nil && highestBidLevel.ToBeDeleted[hbc.ContractID] == nil {
			log.Printf("We do enter here!!")
			if lac.Quantity == hbc.Quantity {
				log.Printf("(1)")
				lowestAskLevel.Orders.Dequeue()
				highestBidLevel.Orders.Dequeue()
				lowestAskLevel.NoOfContracts -= lac.Quantity
				highestBidLevel.NoOfContracts -= hbc.Quantity
				ob.LogHandler(lac, hbc)
			} else if lac.Quantity < hbc.Quantity {
				log.Printf("(2)")
				lowestAskLevel.Orders.Dequeue()
				hbc.Quantity -= lac.Quantity
				lowestAskLevel.NoOfContracts -= lac.Quantity
				ob.LogHandler(lac, hbc)
			} else {
				log.Printf("(3)")
				highestBidLevel.Orders.Dequeue()
				lac.Quantity -= hbc.Quantity
				highestBidLevel.NoOfContracts -= hbc.Quantity
				ob.LogHandler(lac, hbc)
			}
		} else {
			if lowestAskLevel.ToBeDeleted[lac.ContractID] != nil {
				log.Printf("Contract: %s\n Just deleted :O", lac.ContractID)
				lowestAskLevel.Orders.Dequeue()
				lowestAskLevel.NoOfContracts -= lac.Quantity
			}
			if highestBidLevel.ToBeDeleted[hbc.ContractID] != nil {
				log.Printf("Contract: %s\n Just deleted :O", hbc.ContractID)
				highestBidLevel.Orders.Dequeue()
				highestBidLevel.NoOfContracts -= hbc.Quantity
			}
			log.Printf("Checking for a loop lmao")
		}
	}
	log.Printf("Matched %v with %v", lowestAskLevel.Price, highestBidLevel.Price)
}
