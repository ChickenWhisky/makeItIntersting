package subevents

import (
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
)

type SubEvent struct {
	SubEventID     string               // ID used to identify all subevents (within each event there could be multiple YES/NO's) in Ledger as well as used in contracts
	OrderBook      *orderbook.OrderBook // OrderBook used to match contracts for given sub-event
	Yes            bool                 // To identify if it is the Yes or No
	SubEventName   string               // A simple title for the event?
	SubEventExpiry time.Time            // NOT_FINAL
	SubEventStart  time.Time            // NOT_FINAL
	SubEventInfo   string               // NOT_FINAL
	ContractVolume int                  // Metrics for number of contracts issued in the event
	TraderVolume   int                  // Metrics for number of traders in the event
	ValueVolume    int                  // Metrics for number of traders in the event
}

func NewSubEvent(SubEventID string, name string, curTime time.Time, Yes bool) *SubEvent {
	return &SubEvent{
		SubEventID:    SubEventID,
		SubEventName:  name,
		OrderBook:     orderbook.NewOrderBook(),
		SubEventStart: curTime,
		Yes:           Yes,
	}

}
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

func (s *SubEvent) SubmitOrder(o models.Order) {
	s.OrderBook.PushOrderIntoQueue(o)
}
