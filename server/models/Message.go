package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Mime      string `json:"mime"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

func NewMessage() Message {
	return Message{
		Id:        uuid.NewV4().String(),
		Type:      "text",
		Mime:      "text/plain",
		Content:   "",
		CreatedAt: time.Now().Unix(),
	}
}
