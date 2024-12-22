package models

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
)

type OrderQueue struct {
	heap        *priorityqueue.Queue
	orders      map[string]*Contract // Maps IDs to the contract
	toBeDeleted map[string]*Contract // Keeps track of IDs to be deleted
	noOfOrders  int                  // Size of number of orders (Not contracts)
}

// This function should be edited inorder to change the matching formula

// PriceTimeBased implements the basic FIFO match making
func PriceTimeBased(a, b interface{}) int {
	timeA := a.(*Contract).Timestamp
	timeB := b.(*Contract).Timestamp
	priceA := a.(*Contract).Price
	priceB := b.(*Contract).Price
	if priceA == priceB {
		return utils.Int64Comparator(timeA, timeB)
	} else {
		if a.(*Contract).OrderType == "buy" || a.(*Contract).OrderType == "limit_buy" {
			if priceA < priceB {
				return -1
			}
			if priceA > priceB {
				return 1
			}
		} else {
			if priceA > priceB {
				return -1
			}
			if priceA < priceB {
				return 1
			}

		}
	}
	return 0
}

// NewOrderQueue Creates a new instance of the OrderQueue
func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		heap:        priorityqueue.NewWith(PriceTimeBased),
		orders:      make(map[string]*Contract),
		toBeDeleted: make(map[string]*Contract),
		noOfOrders:  0,
	}
}

// Pop returns the contract fromt the top of the OrderQueue as well as dequeues it from the OrderQueue
func (oq OrderQueue) Pop() (*Contract, error) {
	_contract, isPQNotEmpty := oq.heap.Peek()
	if !isPQNotEmpty {
		return nil, errors.New("order queue is empty")
	}
	contract := _contract.(*Contract)
	oq.heap.Dequeue()

	// If the item in contention is to be deleted then simply ignore it and call oq.Pop() again
	if oq.toBeDeleted[contract.ContractID] != nil {
		delete(oq.toBeDeleted, contract.ContractID)
		return oq.Pop()
	}

	// Remove th order from the mapping of contract ID's to Contract pointers as well as reduce the number of orders
	delete(oq.orders, contract.ContractID)
	oq.noOfOrders--

	return contract, nil
}

// Push enqueues a contract into the OrderQueue
func (oq OrderQueue) Push(contract *Contract) error {
	if oq.orders[contract.ContractID] != nil {
		return errors.New("contract already exists")
	}
	oq.orders[contract.ContractID] = contract
	oq.noOfOrders++
	oq.heap.Enqueue(contract)
	return nil
}

// Top returns the contract at the top of the OrderQueue
func (oq OrderQueue) Top() (*Contract, error) {
	_contract, isPQNotEmpty := oq.heap.Peek()
	if !isPQNotEmpty {
		return nil, errors.New("order queue is empty")
	}

	contract := _contract.(*Contract)

	// Checks if the contract needs to be deleted or not
	if oq.toBeDeleted[contract.ContractID] != nil {
		oq.Pop()
		return oq.Top()
	}
	return contract, nil
}

// Delete a contract in the OrderQueue with the contract ID
func (oq OrderQueue) Delete(ID string) error {
	// Implements a lazy deletion sort of method where only the number of orders is reduced now, but it is kept in the
	// toBeDeleted map so that when it appears in the top of the pq it is deleted only then. This is simply to abstract away
	// the deletion of the contract from the queue

	contract := oq.orders[ID]
	if contract == nil {
		return fmt.Errorf("contract %s does not exist", ID)
	}
	if oq.toBeDeleted[contract.ContractID] != nil {
		return fmt.Errorf("contract %s already deleted", ID)
	}

	delete(oq.orders, ID)
	oq.noOfOrders--
	oq.toBeDeleted[ID] = contract
	return nil
}

// Find returns a contract within the OrderQueue if it exists
func (oq OrderQueue) Find(ID string) (*Contract, error) {
	if oq.orders[ID] == nil {
		return nil, fmt.Errorf("contract %s does not exist", ID)
	}
	return oq.orders[ID], nil
}

// clear is a function that simply checks if the contract in contention is to be deleted if the clearing is done the bool
// will be returned as true else it will be returned as false
func (oq OrderQueue) clear(contract *Contract) bool {
	if oq.toBeDeleted[contract.ContractID] != nil {
		delete(oq.toBeDeleted, contract.ContractID)
		return true
	}
	return false
}
