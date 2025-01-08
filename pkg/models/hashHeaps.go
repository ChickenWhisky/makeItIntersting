package models

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
)

type OrderQueue struct {
	heap        *priorityqueue.Queue //  Priority Queue to store the Orders
	orders      map[string]*Order    // Maps IDs to the Order
	toBeDeleted map[string]*Order    // Keeps track of IDs to be deleted
	noOfOrders  int                  // Size of number of orders (Not Orders)
}

// This function should be edited inorder to change the matching formula

// PriceTimeBased implements the basic FIFO match making
func PriceTimeBased(a, b interface{}) int {

	A := a.(*Order)
	B := b.(*Order)
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
		orders:      make(map[string]*Order),
		toBeDeleted: make(map[string]*Order),
		noOfOrders:  0,
	}
}

// Empty checks whether the HashHeap is empty or not
func (oq *OrderQueue) Empty() bool {
	return oq.noOfOrders == 0
}

func (oq *OrderQueue) TopPrice() float32 {
	_Order, err := oq.Top()
	if err != nil {
		log.Printf("Error From TopPrice: %v", err)
		return -1 // Sending -1 as negative pricing doesn't make sense
	}
	return _Order.GetPrice()
}

// Push enqueues a Order into the OrderQueue
func (oq *OrderQueue) Push(Order *Order) error {
	if oq.orders[Order.GetOrderID()] != nil {
		return errors.New("Order already exists")
	}
	oq.orders[Order.GetOrderID()] = Order
	oq.noOfOrders++
	oq.heap.Enqueue(Order)
	return nil
}

// Pop returns the Order from the top of the OrderQueue as well as dequeues it from the OrderQueue
func (oq *OrderQueue) Pop() (*Order, error) {
	_Order, isPQNotEmpty := oq.heap.Dequeue()
	if !isPQNotEmpty {
		return nil, errors.New("order queue is empty")
	}
	Order := _Order.(*Order)

	// If the item in contention is to be deleted then simply ignore it and call oq.Pop() again
	if oq.toBeDeleted[Order.GetOrderID()] == Order {
		delete(oq.toBeDeleted, Order.GetOrderID())
		return oq.Pop()
	}

	// Remove th order from the mapping of Order ID's to Order pointers as well as reduce the number of orders

	delete(oq.orders, Order.GetOrderID())
	oq.noOfOrders--

	return Order, nil
}

// Top returns the Order at the top of the OrderQueue
func (oq *OrderQueue) Top() (*Order, error) {
	_Order, isPQNotEmpty := oq.heap.Peek()
	if !isPQNotEmpty {
		return nil, errors.New("order queue is empty")
	}

	Order := _Order.(*Order)

	// Checks if the Order needs to be deleted or not
	if oq.toBeDeleted[Order.GetOrderID()] == Order {
		oq.Pop()
		return oq.Top()
	}
	return Order, nil
}

// Delete a Order in the OrderQueue with the Order ID
func (oq *OrderQueue) Delete(ID string) error {
	// Implements a lazy deletion sort of method where only the number of orders is reduced now, but it is kept in the
	// toBeDeleted map so that when it appears in the top of the pq it is deleted only then. This is simply to abstract away
	// the deletion of the Order from the queue

	Order := oq.orders[ID]
	if Order == nil {
		return fmt.Errorf("Order %s does not exist", ID)
	}
	if oq.toBeDeleted[Order.GetOrderID()] == Order {
		return fmt.Errorf("Order %s already deleted", ID)
	}
	log.Printf("Order %s is to be deleted FROM ORDERS TRACKER", ID)
	delete(oq.orders, ID)
	oq.noOfOrders--
	log.Printf("Order %s is to added TO BE DELETED TRACKER", ID)
	oq.toBeDeleted[Order.GetOrderID()] = Order
	return nil
}

// Find returns a Order within the OrderQueue if it exists
func (oq *OrderQueue) Find(ID string) (*Order, error) {
	if oq.orders[ID] == nil {
		return nil, fmt.Errorf("Order %s does not exist", ID)
	}
	return oq.orders[ID], nil
}

// clear is a function that simply checks if the Order in contention is to be deleted if the clearing is done the bool
// will be returned as true else it will be returned as false
func (oq *OrderQueue) clear(Order *Order) bool {
	if oq.toBeDeleted[Order.GetOrderID()] == Order {
		delete(oq.toBeDeleted, Order.GetOrderID())
		return true
	}
	return false
}
