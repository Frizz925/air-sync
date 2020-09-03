package repositories

import (
	"air-sync/models"
	"errors"
)

var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrMessageNotFound = errors.New("Message not found")
)

type SessionRepository interface {
	Create() (models.Session, error)
	Find(id string) (models.Session, error)
	InsertMessage(id string, model models.InsertMessage) (models.Message, error)
	DeleteMessage(id string, messageID string) error
	Delete(id string) error
}
