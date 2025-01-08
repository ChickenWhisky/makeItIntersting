package events

import subevents "github.com/ChickenWhisky/makeItIntersting/internals/SubEvents"

// GetEventName returns the name of the event.
func (e *Event) GetEventName() string {
	return e.EventName
}

// SetEventName sets the name of the event.
func (e *Event) SetEventName(en string) {
	e.EventName = en
}

// GetSubEvents returns the sub-events associated with the event.
func (e *Event) GetSubEvents() map[string]*subevents.SubEvent {
	return e.SubEvents
}

// GetEventID returns the ID of the event.
func (e *Event) GetEventID() string {
	return e.EventID
}

// SetEventID sets the ID of the event.
func (e *Event) SetEventID(eventID string) {
	e.EventID = eventID
}

// GetEventInfo returns the information of the event.
func (e *Event) GetEventInfo() string {
	return e.EventInfo
}

// SetEventInfo sets the information of the event.
func (e *Event) SetEventInfo(eventInfo string) {
	e.EventInfo = eventInfo
}

// GetContractVolume returns the contract volume of the event.
func (e *Event) GetContractVolume() int {
	return e.ContractVolume
}

// SetContractVolume sets the contract volume of the event.
func (e *Event) SetContractVolume(contractVolume int) {
	e.ContractVolume = contractVolume
}

// GetTraderVolume returns the trader volume of the event.
func (e *Event) GetTraderVolume() int {
	return e.TraderVolume
}

// SetTraderVolume sets the trader volume of the event.
func (e *Event) SetTraderVolume(traderVolume int) {
	e.TraderVolume = traderVolume
}