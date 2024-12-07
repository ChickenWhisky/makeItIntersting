package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/emirpasic/gods/utils"
)

func LevelByLevel(a, b interface{}) int {
	contractA := a.(*LevelBook)
	contractB := b.(*LevelBook)
	if contractA.Type {
		if contractA.Price < contractB.Price {
			return -1
		}
		if contractA.Price > contractB.Price {
			return 1
		}
		return 0
	} else {
		if contractA.Price > contractB.Price {
			return -1
		}
		if contractA.Price < contractB.Price {
			return 1
		}
		return 0
	}

}

func TimeBased(a, b interface{}) int {
	timeA := a.(models.Contract).Timestamp
	timeB := b.(models.Contract).Timestamp
	return utils.IntComparator(timeA, timeB)
}
