package events

import (
	"time"

	subevents "github.com/ChickenWhisky/makeItIntersting/internals/SubEvents"
)

// GetEventName returns the name of the event.
func (e *Event) GetEventName() string {
	return e.EventName
}

// SetEventName sets the name of the event.
func (e *Event) SetEventName(en string) {
	e.EventName = en
}

// GetSubEvents returns the sub-events associated with the event.
func (e *Event) GetSubEventNames() []string {
	
	names := make([]string, 0, len(e.SubEvents))
	
	for _, subEvent := range e.SubEvents {
		names = append(names,subEvent.GetSubEventName())
	}
	
	return names
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


// GetSubEventSummary returns the summary of the sub-events associated with the event.
func (e *Event) GetSubEventSummary(s string) subevents.SubEventSummary {
	subEventSummary := e.SubEvents[s].Summary()
	return subEventSummary
}

func (e *Event) GetSubEventsSummary() map[string]subevents.SubEventSummary {
	
	subEventsSummary := make(map[string]subevents.SubEventSummary)

	for _, subEvent := range e.SubEvents {
		subEventsSummary[subEvent.GetSubEventID()] = subEvent.Summary()
	}
	return subEventsSummary
	
}

// GetSubEventIDs returns the IDs of the sub-events associated with the event.
func (e *Event) GetSubEventIDs() []string {

	ids := make([]string, 0, len(e.SubEvents))
	
	for _, subEvent := range e.SubEvents {
		ids = append(ids, subEvent.GetSubEventID())
	}
	
	return ids
}

// GetSubEventNameID returns the names and IDs of the sub-events associated with the event.	
func (e *Event) GetSubEventsNameID() [][]string {
	
	var subEventNames []string = e.GetSubEventNames()
	var subEventIDs []string = e.GetSubEventIDs()
	subEventNameID := make([][]string, 2, len(subEventNames))

	for i:=0; i<len(subEventNames);i++{
		subEventNameID = append(subEventNameID, []string{subEventIDs[i], subEventNames[i]})
	}
	
	return subEventNameID
}

// GetSubEventStartTimes returns the start times of the sub-events associated with the event.
func (e *Event) GetSubEventStartTimes() []time.Time {
	
	startTimes := make([]time.Time, 0, len(e.SubEvents))
	
	for _, subEvent := range e.SubEvents {
		startTimes = append(startTimes, subEvent.GetSubEventStart())
	}
	
	return startTimes
}

// GetSubEventEndTimes returns the end times of the sub-events associated with the event.

func (e *Event) GetSubEventEndTimes() []time.Time {
	
	endTimes := make([]time.Time, 0, len(e.SubEvents))
	
	for _, subEvent := range e.SubEvents {
		endTimes = append(endTimes, subEvent.GetSubEventExpiry())
	}
	return endTimes
}