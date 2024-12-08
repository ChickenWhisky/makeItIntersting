package orderbook

import "github.com/ChickenWhisky/makeItIntersting/pkg/models"

// All the contracts that reach this stage of checking exists within the OrderBook

// DeleteContractFromAsks deletes a given contract from the OrderBook
func (ob *OrderBook) DeleteContractFromAsks(contract models.Contract) {
	requiredLevelBook := ob.AsksOrderByOrder[contract.Price]

	// Now update the Level Book
	requiredLevelBook.NoOfContracts -= contract.Quantity
	if requiredLevelBook.NoOfContracts == 0 {
		ob.DeleteLevelBook(requiredLevelBook)
	} else {
		// This adds the system
		delete(requiredLevelBook.Contracts, contract.ContractID)
		requiredLevelBook.ToBeDeleted[contract.ContractID] = &contract
	}
	delete(ob.Orders, contract.ContractID)

}
