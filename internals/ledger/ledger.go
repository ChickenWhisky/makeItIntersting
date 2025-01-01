package ledger

import (
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
)

type Event struct {
	EventID        string               // ID used to identify events in Ledger as well as used in contracts
	OrderBooks     *orderbook.OrderBook // OrderBook used to match contracts for given sub-event
	EventName      string
	EventExpiry    time.Time
	EventStart     time.Time
	EventInfo      string
}

type Ledger struct {
	Events map[string]*Event
}

func newEvent() *Event {
	return &Event{}
}
