package events

type RedisSessionEvent struct {
	SessionEvent
	ClientId string `json:"client_id"`
}
