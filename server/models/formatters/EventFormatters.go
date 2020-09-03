package formatters

import (
	"air-sync/models"
	"air-sync/models/events"
)

func FromSessionEvent(e events.SessionEvent) models.Event {
	evt := models.Event{
		ID:        e.ID,
		Event:     e.Event,
		Data:      e.Value,
		Timestamp: e.Timestamp,
	}
	if e.Error != nil {
		evt.Error = e.Error.Error()
	}
	return evt
}
