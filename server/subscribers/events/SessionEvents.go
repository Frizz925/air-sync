package events

import "air-sync/models"

type SessionUpdate struct {
	Id      string
	Message models.Message
}

type SessionDelete string

const (
	SessionUpdateEventName = "update"
	SessionDeleteEventName = "delete"
)

func SessionEventName(id string) string {
	return "session:" + id
}
