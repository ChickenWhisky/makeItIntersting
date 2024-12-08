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
	delete(ob.Orders, contract.ContractID)

}
