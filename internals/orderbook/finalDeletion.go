package orderbook

import "github.com/ChickenWhisky/makeItIntersting/pkg/models"

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
