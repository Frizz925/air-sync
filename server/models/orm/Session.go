package orm

import (
	"air-sync/models"
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

func ToSessionModel(session Session) models.Session {
	messages := make([]models.Message, len(session.Messages))
	for index, message := range session.Messages {
		messages[index] = ToMessageModel(message)
	}
	return models.Session{
		ID:        session.ID,
		Messages:  messages,
		CreatedAt: session.CreatedAt,
	}
}
