package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Id        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt int64     `json:"created_at"`
}

var EmptySession = Session{}

func NewSession() Session {
	return Session{
		Id:        uuid.NewV4().String(),
		Messages:  make([]Message, 0),
		CreatedAt: time.Now().Unix(),
	}
}
