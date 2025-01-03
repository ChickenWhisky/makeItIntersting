package ledger

import (
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
)




type Ledger struct {
	Events map[string]*Event
}

func newLedger() *Ledger {
	return &Ledger{}
}
