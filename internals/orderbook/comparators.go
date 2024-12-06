package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/emirpasic/gods/utils"
)

func ForAsks(a, b interface{}) int {
	contractA := a.(*models.Contract)
	contractB := b.(*models.Contract)

	if contractA.Price < contractB.Price {
		return -1
	}
	if contractA.Price > contractB.Price {
		return 1
	}

	// If prices are equal, prioritize earlier timestamp
	if contractA.Timestamp < contractB.Timestamp {
		return -1
	}
	if contractA.Timestamp > contractB.Timestamp {
		return 1
	}

	return 0
}
func ForBids(a, b interface{}) int {
	contractA := a.(*models.Contract)
	contractB := b.(*models.Contract)

	if contractA.Price > contractB.Price {
		return -1
	}
	if contractA.Price < contractB.Price {
		return 1
	}

	// If prices are equal, prioritize earlier timestamp
	if contractA.Timestamp < contractB.Timestamp {
		return -1
	}
	if contractA.Timestamp > contractB.Timestamp {
		return 1
	}

	return 0
}
func ForLimitOrdersAsk(a, b interface{}) int {
	priorityA := a.(models.LimitOrderTracker).Price
	priorityB := b.(models.LimitOrderTracker).Price
	return -utils.IntComparator(priorityA, priorityB) // "-" Descending order
}
func ForLimitOrdersBid(a, b interface{}) int {
	priorityA := a.(models.LimitOrderTracker).Price
	priorityB := b.(models.LimitOrderTracker).Price
	return utils.IntComparator(priorityA, priorityB)
}
