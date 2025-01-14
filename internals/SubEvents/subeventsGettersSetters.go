package subevents

import (
	
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/orderbook"
)

type SubEventSummary struct {
	SubEventName   string    // A simple title for the event?
	SubEventStart  time.Time // NOT_FINAL
	SubEventEnd    time.Time // NOT_FINAL
	ContractVolume int       // Metrics for number of contracts issued in the event
	TraderVolume   int       // Metrics for number of traders in the event
	ValueVolume    int       // Metrics for number of traders in the event
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

func (s *SubEvent) GetContractVolume() int {
	return s.ContractVolume
}	

func (s *SubEvent) SetContractVolume(contractVolume int) {
	s.ContractVolume = contractVolume
}

func (s *SubEvent) GetTraderVolume() int {
	return s.TraderVolume
}

func (s *SubEvent) SetTraderVolume(traderVolume int) {
	s.TraderVolume = traderVolume
}

func (s *SubEvent) GetValueVolume() int {
	return s.ValueVolume
}

func (s *SubEvent) SetValueVolume(valueVolume int) {
	s.ValueVolume = valueVolume
}

func (s *SubEvent) Summary() SubEventSummary {
	return SubEventSummary{
		SubEventName:   s.GetSubEventName(),
		SubEventStart:  s.GetSubEventStart(),
		SubEventEnd:    s.GetSubEventExpiry(),
		ContractVolume: s.GetContractVolume(),
		TraderVolume:   s.GetTraderVolume(),
		ValueVolume:    s.GetValueVolume(),
	}
}