package stores

import (
	"air-sync/models"
	"errors"
)

var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrMessageNotFound = errors.New("Message not found")
)

type SessionRepository interface {
	Create() (*models.Session, error)
	Get(id string) (*models.Session, error)
	InsertMessage(id string, message models.Message) error
	DeleteMessage(id string, messageId string) error
	Delete(id string) error
}
