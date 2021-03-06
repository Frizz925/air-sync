package repositories

import (
	"air-sync/models"
	"errors"
	"time"
)

var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrMessageNotFound = errors.New("Message not found")
)

type SessionRepository interface {
	Create() (models.Session, error)
	Find(id string) (models.Session, error)
	FindBefore(t time.Time) ([]models.Session, error)
	InsertMessage(id string, model models.InsertMessage) (models.Message, error)
	DeleteMessage(id string, messageId string) error
	Delete(id string) error
	DeleteMany(ids []string) (int, error)
}
