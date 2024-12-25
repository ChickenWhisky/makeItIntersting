package orderbook

//
//import (
//	"github.com/ChickenWhisky/makeItIntersting/pkg/helpers"
//	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
//	"github.com/emirpasic/gods/queues/priorityqueue"
//	"os"
//)
//
//func (ob *OrderBook) AddContractToAsks(contract models.Contract) {
//
//	// Extract pointer to the required level
//	requiredLevel, existsInOrderBook := ob.AsksOrderByOrder[contract.Price]
//	if existsInOrderBook {
//		requiredLevel.NoOfContracts += contract.Quantity
//		requiredLevel.Orders.Enqueue(&contract)
//		requiredLevel.Contracts[contract.ContractID] = &contract
//	} else {
//
//		newLevel := &LevelBook{
//			LevelID:       helpers.GenerateRandomString(helpers.ConvertStringToInt(os.Getenv("LEVEL_ID_LENGTH"))),
//			Price:         contract.Price,
//			Type:          true,
//			NoOfContracts: contract.Quantity,
//			Orders:        priorityqueue.NewWith(TimeBased),
//			ToBeDeleted:   make(map[string]*models.Contract),
//			Contracts:     make(map[string]*models.Contract),
//		}
//		newLevel.Contracts[contract.ContractID] = &contract
//		newLevel.Orders.Enqueue(&contract)
//		ob.AsksOrderByOrder[contract.Price] = newLevel
//		ob.AsksLevelByLevel.Enqueue(newLevel)
//	}
//
//}
//func (ob *OrderBook) AddContractToBids(contract models.Contract) {
//
//	// Extract pointer to the required level
//	requiredLevel, existsInOrderBook := ob.BidsOrderByOrder[contract.Price]
//	if existsInOrderBook {
//		requiredLevel.NoOfContracts += contract.Quantity
//		requiredLevel.Orders.Enqueue(&contract)
//		requiredLevel.Contracts[contract.ContractID] = &contract
//	} else {
//
//		newLevel := &LevelBook{
//			LevelID:       helpers.GenerateRandomString(helpers.ConvertStringToInt(os.Getenv("LEVEL_ID_LENGTH"))),
//			Price:         contract.Price,
//			Type:          true,
//			NoOfContracts: contract.Quantity,
//			Orders:        priorityqueue.NewWith(TimeBased),
//			ToBeDeleted:   make(map[string]*models.Contract),
//			Contracts:     make(map[string]*models.Contract),
//		}
//		newLevel.Contracts[contract.ContractID] = &contract
//		newLevel.Orders.Enqueue(&contract)
//		ob.BidsOrderByOrder[contract.Price] = newLevel
//		ob.BidsLevelByLevel.Enqueue(newLevel)
//	}
//
//}
//func (ob *OrderBook) AddContractToLimitAsks(contract models.Contract) {
//	// Extract pointer to the required level
//	requiredLevel, existsInOrderBook := ob.LimitAsksOrderByOrder[contract.Price]
//	if existsInOrderBook {
//		requiredLevel.NoOfContracts += contract.Quantity
//		requiredLevel.Orders.Enqueue(&contract)
//		requiredLevel.Contracts[contract.ContractID] = &contract
//	} else {
//
//		newLevel := &LevelBook{
//			LevelID:       helpers.GenerateRandomString(helpers.ConvertStringToInt(os.Getenv("LEVEL_ID_LENGTH"))),
//			Price:         contract.Price,
//			Type:          true,
//			NoOfContracts: contract.Quantity,
//			Orders:        priorityqueue.NewWith(TimeBased),
//			ToBeDeleted:   make(map[string]*models.Contract),
//			Contracts:     make(map[string]*models.Contract),
//		}
//		newLevel.Contracts[contract.ContractID] = &contract
//		newLevel.Orders.Enqueue(&contract)
//		ob.LimitAsksOrderByOrder[contract.Price] = newLevel
//		ob.LimitAsksLevelByLevel.Enqueue(newLevel)
//	}
//}
//func (ob *OrderBook) AddContractToLimitBids(contract models.Contract) {
//	// Extract pointer to the required level
//	requiredLevel, existsInOrderBook := ob.LimitBidsOrderByOrder[contract.Price]
//	if existsInOrderBook {
//		requiredLevel.NoOfContracts += contract.Quantity
//		requiredLevel.Orders.Enqueue(&contract)
//		requiredLevel.Contracts[contract.ContractID] = &contract
//	} else {
//
//		newLevel := &LevelBook{
//			LevelID:       helpers.GenerateRandomString(helpers.ConvertStringToInt(os.Getenv("LEVEL_ID_LENGTH"))),
//			Price:         contract.Price,
//			Type:          true,
//			NoOfContracts: contract.Quantity,
//			Orders:        priorityqueue.NewWith(TimeBased),
//			ToBeDeleted:   make(map[string]*models.Contract),
//			Contracts:     make(map[string]*models.Contract),
//		}
//		newLevel.Contracts[contract.ContractID] = &contract
//		newLevel.Orders.Enqueue(&contract)
//		ob.LimitBidsOrderByOrder[contract.Price] = newLevel
//		ob.LimitBidsLevelByLevel.Enqueue(newLevel)
//	}
//}
