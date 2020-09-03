package events

import (
	"air-sync/models"

	uuid "github.com/satori/go.uuid"
)

type BaseEvent struct {
	ID        string      `json:"id"`
	Event     string      `json:"event"`
	Value     interface{} `json:"data,omitempty"`
	Error     error       `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

type Event BaseEvent

type SessionEvent struct {
	BaseEvent
	SessionID string `json:"session_id"`
}

type SessionCreate models.Session

type SessionDelete string

type MessageInsert struct {
	SessionID string         `json:"session_id"`
	Message   models.Message `json:"message"`
}

type MessageDelete struct {
	SessionID string `json:"session_id"`
	MessageID string `json:"message_id"`
}

const (
	EventSession         = "session"
	EventSessionCreated  = "session.created"
	EventSessionDeleted  = "session.deleted"
	EventMessageInserted = "message.inserted"
	EventMessageDeleted  = "message.deleted"
)

func CreateEvent(event string, v interface{}, err error) Event {
	return Event{
		ID:        uuid.NewV4().String(),
		Event:     event,
		Value:     v,
		Error:     err,
		Timestamp: models.Timestamp(),
	}
}

func CreateSessionEvent(id string, event string, v interface{}, err error) SessionEvent {
	return SessionEvent{
		BaseEvent: BaseEvent(CreateEvent(event, v, err)),
		SessionID: id,
	}
}

func EventSessionID(id string) string {
	return "session:" + id
}
