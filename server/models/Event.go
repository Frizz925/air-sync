package models

type Event struct {
	ID        string      `json:"id"`
	Event     string      `json:"event"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}
