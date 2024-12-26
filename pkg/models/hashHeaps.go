package models

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
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

	A := a.(*Contract)
	B := b.(*Contract)
	if A.GetPrice() == B.GetPrice() {
		return utils.Int64Comparator(A.GetTimestamp(), B.GetTimestamp())
	} else {
		if A.GetOrderType() == "buy" || A.GetOrderType() == "limit_buy" {
			if A.GetPrice() < B.GetPrice() {
				return 1
			}
			if A.GetPrice() > B.GetPrice() {
				return -1
			}
		} else {
			if A.GetPrice() > B.GetPrice() {
				return 1
			}
			if A.GetPrice() < B.GetPrice() {
				return -1
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

// Empty checks whether the HashHeap is empty or not
func (oq *OrderQueue) Empty() bool {
	return oq.noOfOrders == 0
}

func (oq *OrderQueue) TopPrice() float32 {
	_contract, err := oq.Top()
	if err != nil {
		log.Printf("Error From TopPrice: %v", err)
		return -1 // Sending -1 as negative pricing doesn't make sense
	}
	return _contract.GetPrice()
}

// Push enqueues a contract into the OrderQueue
func (oq *OrderQueue) Push(contract *Contract) error {
	if oq.orders[contract.GetContractID()] != nil {
		return errors.New("contract already exists")
	}
	oq.orders[contract.GetContractID()] = contract
	oq.noOfOrders++
	oq.heap.Enqueue(contract)
	return nil
}

// Pop returns the contract from the top of the OrderQueue as well as dequeues it from the OrderQueue
func (oq *OrderQueue) Pop() (*Contract, error) {
	_contract, isPQNotEmpty := oq.heap.Dequeue()
	if !isPQNotEmpty {
		return nil, errors.New("order queue is empty")
	}
	contract := _contract.(*Contract)

	// If the item in contention is to be deleted then simply ignore it and call oq.Pop() again
	if oq.toBeDeleted[contract.GetContractID()] == contract {
		delete(oq.toBeDeleted, contract.GetContractID())
		return oq.Pop()
	}

	// Remove th order from the mapping of contract ID's to Contract pointers as well as reduce the number of orders

	delete(oq.orders, contract.GetContractID())
	oq.noOfOrders--

	return contract, nil
}

// Top returns the contract at the top of the OrderQueue
func (oq *OrderQueue) Top() (*Contract, error) {
	_contract, isPQNotEmpty := oq.heap.Peek()
	if !isPQNotEmpty {
		return nil, errors.New("order queue is empty")
	}

	contract := _contract.(*Contract)

	// Checks if the contract needs to be deleted or not
	if oq.toBeDeleted[contract.GetContractID()] == contract {
		oq.Pop()
		return oq.Top()
	}
	return contract, nil
}

// Delete a contract in the OrderQueue with the contract ID
func (oq *OrderQueue) Delete(ID string) error {
	// Implements a lazy deletion sort of method where only the number of orders is reduced now, but it is kept in the
	// toBeDeleted map so that when it appears in the top of the pq it is deleted only then. This is simply to abstract away
	// the deletion of the contract from the queue

	contract := oq.orders[ID]
	if contract == nil {
		return fmt.Errorf("contract %s does not exist", ID)
	}
	if oq.toBeDeleted[contract.GetContractID()] == contract {
		return fmt.Errorf("contract %s already deleted", ID)
	}
	log.Printf("Contract %s is to be deleted FROM ORDERS TRACKER", ID)
	delete(oq.orders, ID)
	oq.noOfOrders--
	log.Printf("Contract %s is to added TO BE DELETED TRACKER", ID)
	oq.toBeDeleted[contract.GetContractID()] = contract
	return nil
}

// Find returns a contract within the OrderQueue if it exists
func (oq *OrderQueue) Find(ID string) (*Contract, error) {
	if oq.orders[ID] == nil {
		return nil, fmt.Errorf("contract %s does not exist", ID)
	}
	return oq.orders[ID], nil
}

// clear is a function that simply checks if the contract in contention is to be deleted if the clearing is done the bool
// will be returned as true else it will be returned as false
func (oq *OrderQueue) clear(contract *Contract) bool {
	if oq.toBeDeleted[contract.GetContractID()] == contract {
		delete(oq.toBeDeleted, contract.GetContractID())
		return true
	}
	return false
}
