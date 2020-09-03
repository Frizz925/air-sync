package models

import (
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt int64     `json:"created_at"`
}

var EmptySession = Session{}

func NewSession() Session {
	return Session{
		ID:        uuid.NewV4().String(),
		Messages:  make([]Message, 0),
		CreatedAt: Timestamp(),
	}
}
