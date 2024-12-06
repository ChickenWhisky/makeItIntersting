package models

import (
	"github.com/emirpasic/gods/queues/priorityqueue"
	"sync"
)

// OrderBook stores order data and handles order processing.

type OrderBook struct {
	Asks              *priorityqueue.Queue
	Bids              *priorityqueue.Queue
	LimitOrderAsks    *priorityqueue.Queue
	LimitOrderBids    *priorityqueue.Queue
	IncomingContracts chan Contract
	UserOrders        map[string][]Contract
	LastMatchedPrices []float64
	mu                sync.Mutex
}
