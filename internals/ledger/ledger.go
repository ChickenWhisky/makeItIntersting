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

func NewLedger() *Ledger {
	return &Ledger{
		Events: make(map[string]*events.Event),
	}
}

func (l *Ledger) addEvent(e string, subEvents []string) {
	// Create a new event
	_, err := l.Events[e]
	if err {
		log.Printf("Event already exists")
		return

	}

	event, errr := events.NewEvents(e, subEvents)
	if errr != nil {
		log.Printf("Error in adding event : %v", err)
	}
	l.Events[e] = event

}

func (l *Ledger) SubmitOrder(o models.Order) error {
	// Check if the event exists
	if l.Events[o.GetEventID()] == nil {
		log.Printf("Event doesn't exist")
		return errors.New("Event doesn't exist")
	}
	
	// Submit the order to the respective event
	event,eventExists := l.Events[o.GetEventID()]
	if !eventExists{
		log.Printf("Error in submitting order ")
		return errors.New("Error in submitting order")
	}
	err := event.SubmitOrder(o)
	if err != nil {
		log.Printf("Error in submitting order : %v", err)
		return err
	}
	log.Printf("Order submitted successfully")
	return nil
}


type EventSummary struct {
	EventID   string
	SubEvents []string
}

func (l *Ledger) GetEvent(e string) (EventSummary, error) {
	event, eventExists := l.Events[e]
	if !eventExists {
		log.Printf("Event doesn't exist")
		return EventSummary{}, errors.New("Event doesn't exist")
	}
	var summary EventSummary
	summary.EventID = e
	summary.SubEvents = event.GetSubEventNames()
	return summary, nil
}



func (l *Ledger) GetEvents() map[string]*EventSummary {
	// returns only event_id and a list of subevents in a struct
	
	summaries := make(map[string]*EventSummary)
	for id, event := range l.Events {
		summaries[id] = &EventSummary{
			EventID:   id,
			SubEvents: event.GetSubEventNames(),
		}
	}
	return summaries


}