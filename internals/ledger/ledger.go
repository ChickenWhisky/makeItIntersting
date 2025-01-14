package ledger

import (
	"errors"
	"log"


	"github.com/ChickenWhisky/makeItIntersting/internals/events"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
)

type Ledger struct {
	Events map[string]*events.Event
}

// NewLedger creates a new instance of Ledger
func NewLedger() *Ledger {

	return &Ledger{
		Events: make(map[string]*events.Event),
	}
}

// AddEvent adds a new event to the ledger
func (l *Ledger) AddEvent(e string, subEvents []string) error {

	// Check if the event already exists
	_, exists := l.Events[e]
	if !exists {
		log.Printf("Event already exists")
		return errors.New("Event already exists")
	}

	// Create a new event
	event, err := events.NewEvents(e, subEvents)
	if err != nil {
		log.Printf("Error in adding event : %v", err)
		return err
	}

	// Add event to the map of events in Ledger
	l.Events[e] = event
	log.Printf("Event added successfully : %v", e)
	return nil
}

// SubmitOrder submits an order to the respective event
func (l *Ledger) SubmitOrder(o models.Order) error {

	// Check if the event exists
	event, eventExists := l.Events[o.GetEventID()]
	if !eventExists {
		log.Printf("Error in submitting order ")
		return errors.New("Event doesn't exist")
	}

	// Submit the order to the event
	err := event.SubmitOrder(o)
	if err != nil {
		log.Printf("Error in submitting order : %v", err)
		return err
	}

	log.Printf("Order submitted successfully")
	return nil
}



func (l *Ledger) GetEvent(e string) (events.EventSummary, error) {
	event, eventExists := l.Events[e]
	if !eventExists {
		log.Printf("Event doesn't exist")
		return events.EventSummary{}, errors.New("Event doesn't exist")
	}
	var summary events.EventSummary
	summary.EventID = e
	summary.EventName = event.EventName
	summary.SubEvents = event.GetSubEventsSummary()
	
	return summary, nil
}

func (l *Ledger) GetEvents() []events.EventSummary {
	// returns only event_id and a list of subevents in a struct

	eventList := make([]events.EventSummary, 0, len(l.Events))
	for _, event := range l.Events {
		var summary events.EventSummary
		summary.EventID = event.GetEventID()
		summary.EventName = event.GetEventName()
		summary.SubEvents = event.GetSubEventsSummary()
		eventList = append(eventList, summary)
	}	

	return eventList
}
