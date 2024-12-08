package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/emirpasic/gods/queues/priorityqueue"
)

func (ob *OrderBook) AddContractToAsks(contract models.Contract) {

	// Extract pointer to the required level
	requiredLevel, existsInOrderBook := ob.AsksOrderByOrder[contract.Price]
	if existsInOrderBook {
		requiredLevel.NoOfContracts += contract.Quantity
		requiredLevel.Orders.Enqueue(contract)
		requiredLevel.Contracts[contract.ContractID] = &contract
	} else {

		newLevel := &LevelBook{
			Price:         contract.Price,
			Type:          true,
			NoOfContracts: contract.Quantity,
			Orders:        priorityqueue.NewWith(TimeBased),
			ToBeDeleted:   make(map[string]*models.Contract),
			Contracts:     make(map[string]*models.Contract),
		}
		newLevel.Orders.Enqueue(contract)
		ob.AsksLevelByLevel.Enqueue(newLevel)
	}

}
func (ob *OrderBook) AddContractToBids(contract models.Contract) {

	// Extract pointer to the required level
	requiredLevel, existsInOrderBook := ob.BidsOrderByOrder[contract.Price]
	if existsInOrderBook {
		requiredLevel.NoOfContracts += contract.Quantity
		requiredLevel.Orders.Enqueue(contract)
		requiredLevel.Contracts[contract.ContractID] = &contract
	} else {

		newLevel := &LevelBook{
			Price:         contract.Price,
			Type:          true,
			NoOfContracts: contract.Quantity,
			Orders:        priorityqueue.NewWith(TimeBased),
			ToBeDeleted:   make(map[string]*models.Contract),
			Contracts:     make(map[string]*models.Contract),
		}
		newLevel.Orders.Enqueue(contract)
		ob.BidsLevelByLevel.Enqueue(newLevel)
	}

}
func (ob *OrderBook) AddContractToLimitAsks(contract models.Contract) {
	// Extract pointer to the required level
	requiredLevel, existsInOrderBook := ob.LimitAsksOrderByOrder[contract.Price]
	if existsInOrderBook {
		requiredLevel.NoOfContracts += contract.Quantity
		requiredLevel.Orders.Enqueue(contract)
		requiredLevel.Contracts[contract.ContractID] = &contract
	} else {

		newLevel := &LevelBook{
			Price:         contract.Price,
			Type:          true,
			NoOfContracts: contract.Quantity,
			Orders:        priorityqueue.NewWith(TimeBased),
			ToBeDeleted:   make(map[string]*models.Contract),
			Contracts:     make(map[string]*models.Contract),
		}
		newLevel.Orders.Enqueue(contract)
		ob.LimitAsksLevelByLevel.Enqueue(newLevel)
	}
}
func (ob *OrderBook) AddContractToLimitBids(contract models.Contract) {
	// Extract pointer to the required level
	requiredLevel, existsInOrderBook := ob.LimitBidsOrderByOrder[contract.Price]
	if existsInOrderBook {
		requiredLevel.NoOfContracts += contract.Quantity
		requiredLevel.Orders.Enqueue(contract)
		requiredLevel.Contracts[contract.ContractID] = &contract
	} else {

		newLevel := &LevelBook{
			Price:         contract.Price,
			Type:          true,
			NoOfContracts: contract.Quantity,
			Orders:        priorityqueue.NewWith(TimeBased),
			ToBeDeleted:   make(map[string]*models.Contract),
			Contracts:     make(map[string]*models.Contract),
		}
		newLevel.Orders.Enqueue(contract)
		ob.LimitBidsLevelByLevel.Enqueue(newLevel)
	}
}