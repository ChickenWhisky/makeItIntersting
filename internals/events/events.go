package events

import (
	"errors"
	"log"
	"time"

	"github.com/ChickenWhisky/makeItIntersting/internals/SubEvents"
	"github.com/ChickenWhisky/makeItIntersting/pkg/helpers"
)

type Event struct {
	EventID        string
	EventInfo      string
	EventName      string
	SubEvents      map[string]*subevents.SubEvent // map it based on SubEventID
	ContractVolume int                            // Metrics for number of contracts issued in the event
	TraderVolume   int                            // Metrics for number of traders in the event
}

func NewEvents(en string, SubEventNames []string) (*Event, error) {
	m := make(map[string]*subevents.SubEvent)
	t := time.Now()
	if len(SubEventNames)%2 != 0 {
		return nil, errors.New("every SubEvent should have a corresponding YES/NO SubEvent")
	}

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
		EventName:      en,
		SubEvents:      m,
		ContractVolume: 0,
		TraderVolume:   0,
	}, nil
}
