package events

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/SubEvents"
	"github.com/ChickenWhisky/makeItIntersting/pkg/helpers"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
)

type Event struct {
	EventID        string
	EventInfo      string
	EventName      string
	SubEvents      map[string]*subevents.SubEvent // map it based on SubEventID
	ContractVolume int                            // Metrics for number of contracts issued in the event
	TraderVolume   int                            // Metrics for number of traders in the event
}
type EventSummary struct {
	EventID   string                               // ID used to identify all events in Ledger as well as used in contracts
	EventName string                               // A simple title for the event?
	SubEvents map[string]subevents.SubEventSummary // A map of subevents in the event
}

func NewEvents(en string, SubEventNames []string) (*Event, error) {

	// Create a map to store SubEvents
	m := make(map[string]*subevents.SubEvent)
	t := time.Now()

	if len(SubEventNames)%2 != 0 {
		return nil, errors.New("every SubEvent should have a corresponding YES/NO SubEvent")
	}

	// Create SubEvents for the Event and validate names
	for _, s := range SubEventNames {
		err := helpers.ValidateSubEventName(s)
		if err == nil {
			yes := false
			if s[len(s)-3:] == "_YES" {
				yes = true
			}
			hashedID := helpers.HashText(s)
			m[hashedID] = subevents.NewSubEvent(hashedID, s, t, yes)
		} else {
			log.Printf("Error in creating SubEvent : %v", err)
		}

	}
	return &Event{
		EventID:        helpers.HashText(en),
		EventName:      en,
		SubEvents:      m,
		ContractVolume: 0,
		TraderVolume:   0,
	}, nil
}

// AddSubevent adds subevents to the event
func (e *Event) AddSubevent(sl []subevents.SubEvent) error {
	m := make(map[string]*subevents.SubEvent)
	
	if len(sl)%2 != 0 {
		return errors.New("every SubEvent should have a corresponding YES/NO SubEvent")
	}

	// Create SubEvents for the Event and validate names
	for _, s := range sl {
		err := helpers.ValidateSubEventName(s.SubEventName)
		if err == nil {
			yes := false
			if s.GetSubEventName()[len(s.GetSubEventName())-3:] == "_YES" {
				yes = true
			}
			s.SetYes(yes)
			s.SubEventID = helpers.HashText(s.SubEventName)
			m[s.GetSubEventID()] = &s 
		} else {
			log.Printf("Error in creating SubEvent : %v", err)
		}

	}
	return nil
}
// SubmitOrder submits an order to the respective SubEvent
func (e *Event) SubmitOrder(o models.Order) error {

	// Check if the SubEvent exists
	_, ok := e.SubEvents[o.GetSubEventID()]
	if !ok {
		fmt.Print(e.SubEvents[o.GetSubEventID()])
		log.Printf("SubEvent doesn't exist")
		return errors.New("SubEvent doesn't exist")
	}

	// Submit the order to the respective SubEvent
	err := e.SubEvents[o.GetSubEventID()].SubmitOrder(o)
	if err != nil {
		log.Printf("Error in submitting order : %v", err)
		return err
	}

	return nil
}
