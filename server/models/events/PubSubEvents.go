package events

type PubSubSessionEvent struct {
	SessionEvent
	ClientID string `json:"client_id"`
}
