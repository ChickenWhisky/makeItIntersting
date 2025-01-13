package subevents

import (
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
)

func (s *SubEvent) GetSubEventID() string {
	return s.SubEventID
}

func (s *SubEvent) SetSubEventID(id string) {
	s.SubEventID = id
}

func (s *SubEvent) GetOrderBook() *orderbook.OrderBook {
	return s.OrderBook
}

func (s *SubEvent) SetOrderBook(ob *orderbook.OrderBook) {
	s.OrderBook = ob
}

func (s *SubEvent) GetYes() bool {
	return s.Yes
}

func (s *SubEvent) SetYes(yes bool) {
	s.Yes = yes
}

func (s *SubEvent) GetSubEventName() string {
	return s.SubEventName
}

func (s *SubEvent) SetSubEventName(name string) {
	s.SubEventName = name
}

func (s *SubEvent) GetSubEventExpiry() time.Time {
	return s.SubEventExpiry
}

func (s *SubEvent) SetSubEventExpiry(expiry time.Time) {
	s.SubEventExpiry = expiry
}

func (s *SubEvent) GetSubEventStart() time.Time {
	return s.SubEventStart
}

func (s *SubEvent) SetSubEventStart(start time.Time) {
	s.SubEventStart = start
}

func (s *SubEvent) GetSubEventInfo() string {
	return s.SubEventInfo
}

func (s *SubEvent) SetSubEventInfo(info string) {
	s.SubEventInfo = info
}
