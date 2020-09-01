package events

import "air-sync/models"

type SessionEvent struct {
	SessionId string      `json:"session_id"`
	Event     string      `json:"event"`
	Value     interface{} `json:"data,omitempty"`
	Error     error       `json:"error,omitempty"`
}

type SessionCreate models.Session

type SessionDelete string

type MessageInsert struct {
	SessionId string         `json:"session_id"`
	Message   models.Message `json:"message"`
}

type MessageDelete struct {
	SessionId string `json:"session_id"`
	MessageId string `json:"message_id"`
}

const (
	EventSession         = "session"
	EventSessionCreated  = "session.created"
	EventSessionDeleted  = "session.deleted"
	EventMessageInserted = "message.inserted"
	EventMessageDeleted  = "message.deleted"
)

func EventSessionId(id string) string {
	return "session:" + id
}
