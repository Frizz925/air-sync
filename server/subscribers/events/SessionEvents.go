package events

import "air-sync/models"

type SessionEvent struct {
	Event string
	Value interface{}
	Error error
}

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
	SessionDeleted  = "session/delete"
	MessageInserted = "message/insert"
	MessageDeleted  = "message/delete"
)

func SessionEventName(id string) string {
	return "session:" + id
}
