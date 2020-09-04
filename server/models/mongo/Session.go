package mongo

import (
	"air-sync/models"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID        string `bson:"id"`
	CreatedAt int64  `bson:"created_at"`
}

func NewSession() Session {
	return Session{
		ID:        uuid.NewV4().String(),
		CreatedAt: models.Timestamp(),
	}
}

func ToSessionModel(session Session, messages []models.Message) models.Session {
	return models.Session{
		ID:        session.ID,
		Messages:  messages,
		CreatedAt: session.CreatedAt,
	}
}
