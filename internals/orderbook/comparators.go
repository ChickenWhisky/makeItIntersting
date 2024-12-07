package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/emirpasic/gods/utils"
)

func ForAsksLevelByLevel(a, b interface{}) int {
	contractA := a.(*LevelBook)
	contractB := b.(*LevelBook)

	if contractA.Price < contractB.Price {
		return -1
	}
	if contractA.Price > contractB.Price {
		return 1
	}
	return 0
}
func ForBidsLevelByLevel(a, b interface{}) int {
	contractA := a.(*LevelBook)
	contractB := b.(*LevelBook)

	if contractA.Price > contractB.Price {
		return -1
	}
	if contractA.Price < contractB.Price {
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

func TimeBased(a, b interface{}) int {
	timeA := a.(models.Contract).Timestamp
	timeB := b.(models.Contract).Timestamp
	return utils.IntComparator(timeA, timeB)
}
