package repositories

import (
	"air-sync/models"
	"errors"
	"time"
)

var ErrAttachmentNotFound = errors.New("Attachment not found")

type AttachmentRepository interface {
	Create(arg models.CreateAttachment) (models.Attachment, error)
	Find(id string) (models.Attachment, error)
	FindOrphansBefore(t time.Time) ([]models.Attachment, error)
	Delete(id string) error
	DeleteMany(ids []string) (int, error)
}
