package formatters

import (
	"air-sync/models"
	"air-sync/models/events"
)

func FromSessionEvent(e events.SessionEvent) models.Event {
	evt := models.Event{
		Id:        e.Id,
		Event:     e.Event,
		Data:      e.Value,
		Timestamp: e.Timestamp,
	}
	if e.Error != nil {
		evt.Error = e.Error.Error()
	}
	return evt
}
