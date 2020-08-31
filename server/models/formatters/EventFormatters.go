package formatters

import (
	"air-sync/models"
	"air-sync/subscribers/events"
	"air-sync/util"
)

func FromSessionEvent(e events.SessionEvent) models.Event {
	evt := models.Event{
		Event:     e.Event,
		Data:      e.Value,
		Timestamp: util.TimeNow(),
	}
	if e.Error != nil {
		evt.Error = e.Error.Error()
	}
	return evt
}
