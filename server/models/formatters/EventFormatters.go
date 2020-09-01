package formatters

import (
	"air-sync/models"
	"air-sync/subscribers/events"
	"time"
)

func FromSessionEvent(e events.SessionEvent) models.Event {
	evt := models.Event{
		Event:     e.Event,
		Data:      e.Value,
		Timestamp: time.Now().Unix(),
	}
	if e.Error != nil {
		evt.Error = e.Error.Error()
	}
	return evt
}
