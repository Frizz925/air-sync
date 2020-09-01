package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID        string `gorm:"primaryKey"`
	Messages  []Message
	CreatedAt int64 `gorm:"autoCreateTime"`
}

func NewSession() Session {
	return Session{
		ID:        uuid.NewV4().String(),
		Messages:  make([]Message, 0),
		CreatedAt: time.Now().Unix(),
	}
}
