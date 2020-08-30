package stores

import (
	"air-sync/models"
	"errors"
)

var ErrSessionNotFound = errors.New("Session not found")

type SessionRepository interface {
	Create() (*models.Session, error)
	Get(id string) (*models.Session, error)
	Update(id string, message models.Message) error
	Delete(id string) error
}
