package repositories

import (
	"air-sync/models"
	"air-sync/repositories/entities"
	"errors"
)

var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrMessageNotFound = errors.New("Message not found")
)

type SessionRepository interface {
	Create() (entities.Session, error)
	Find(id string) (entities.Session, error)
	InsertMessage(id string, model models.Message) (entities.Message, error)
	DeleteMessage(id string, messageId string) error
	Delete(id string) error
}
