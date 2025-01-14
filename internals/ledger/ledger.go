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

func (l *Ledger) AddEvent(e string, subEvents []string) {
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
	EventName string
	SubEvents [][]string
}

func (l *Ledger) GetEvent(e string) (EventSummary, error) {
	event, eventExists := l.Events[e]
	if !eventExists {
		log.Printf("Event doesn't exist")
		return EventSummary{}, errors.New("Event doesn't exist")
	}
	var summary EventSummary
	summary.EventID = e
	summary.EventName = event.EventName

	// Subevents is a 2D array in which each row contains the subevent ID and the name of the subevent
	var subEventNames []string = event.GetSubEventNames()
	var subEventIDs []string = event.GetSubEventIDs()

	summary.SubEvents = make([][]string, 2, len(subEventNames))
	for i:=0; i<len(subEventNames);i++{
		summary.SubEvents = append(summary.SubEvents, []string{subEventIDs[i], subEventNames[i]})
	}

	return summary, nil
}



func (l *Ledger) GetEvents() []EventSummary {
	// returns only event_id and a list of subevents in a struct
	
	var summaries []EventSummary
	for id, event := range l.Events {
		summaries = append(summaries, EventSummary{
			EventID:   id,
			EventName: event.EventName,
			SubEvents: event.GetSubEventsNameID(),
		})
	}
	return summaries


}