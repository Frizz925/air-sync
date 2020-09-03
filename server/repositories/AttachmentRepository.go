package repositories

import (
	"air-sync/models"
	"errors"
)

var ErrAttachmentNotFound = errors.New("Attachment not found")

type AttachmentRepository interface {
	Create(arg models.CreateAttachment) (models.Attachment, error)
	Find(id string) (models.Attachment, error)
	Delete(id string) error
}
