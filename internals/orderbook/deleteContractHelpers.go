package orderbook

import "github.com/ChickenWhisky/makeItIntersting/pkg/models"

// All the contracts that reach this stage of checking exists within the OrderBook

// DeleteContractFromAsks deletes a given contract from the OrderBook
func (ob *OrderBook) DeleteContractFromAsks(contract models.Contract) {
	requiredLevelBook := ob.AsksOrderByOrder[contract.Price]

	// Now update the Level Book

	if requiredLevelBook.NoOfContracts == contract.Quantity {
		ob.DeleteLevelBook(requiredLevelBook)
	} else {
		// This adds the system
		// We don't remove it from the main list of contracts as we
		// can only do that once the respective contract pops up at the top
		// of the heap
		requiredLevelBook.ToBeDeleted[contract.ContractID] = &contract
	}
}

// DeleteContractFromBids deletes a given contract from the OrderBook
func (ob *OrderBook) DeleteContractFromBids(contract models.Contract) {
	requiredLevelBook := ob.BidsOrderByOrder[contract.Price]

	// Now update the Level Book

	if requiredLevelBook.NoOfContracts == contract.Quantity {
		ob.DeleteLevelBook(requiredLevelBook)
	} else {
		// This adds the system
		// We don't remove it from the main list of contracts as we
		// can only do that once the respective contract pops up at the top
		// of the heap
		requiredLevelBook.ToBeDeleted[contract.ContractID] = &contract
	}
}

// DeleteContractFromLimitAsks deletes a given contract from the OrderBook
func (ob *OrderBook) DeleteContractFromLimitAsks(contract models.Contract) {
	requiredLevelBook := ob.LimitAsksOrderByOrder[contract.Price]

	// Now update the Level Book

	if requiredLevelBook.NoOfContracts == contract.Quantity {
		ob.DeleteLevelBook(requiredLevelBook)
	} else {
		// This adds the system
		// We don't remove it from the main list of contracts as we
		// can only do that once the respective contract pops up at the top
		// of the heap
		requiredLevelBook.ToBeDeleted[contract.ContractID] = &contract
	}
}

// DeleteContractFromLimitBids deletes a given contract from the OrderBook
func (ob *OrderBook) DeleteContractFromLimitBids(contract models.Contract) {
	requiredLevelBook := ob.LimitBidsOrderByOrder[contract.Price]

	// Now update the Level Book

	if requiredLevelBook.NoOfContracts == contract.Quantity {
		ob.DeleteLevelBook(requiredLevelBook)
	} else {
		// This adds the system
		// We don't remove it from the main list of contracts as we
		// can only do that once the respective contract pops up at the top
		// of the heap
		requiredLevelBook.ToBeDeleted[contract.ContractID] = &contract
	}
}
