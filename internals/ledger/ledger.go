package ledger

import (
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/events"
	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
)




type Ledger struct {
	Events map[string]*events.Event
}

func newLedger() *Ledger {
	return &Ledger{}
}
