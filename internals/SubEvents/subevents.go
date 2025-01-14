package subevents

import (
	"errors"
	"log"
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

func (s *SubEvent) SubmitOrder(o models.Order) error {
	if time.Now().After(s.SubEventExpiry) {
		log.Printf("SubEvent has expired")
		return errors.New("SubEvent has expired")
	}
	s.OrderBook.PushOrderIntoQueue(o)
	return nil
 }
