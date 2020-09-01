package repositories

import (
	"air-sync/repositories/entities"
	"errors"
)

var ErrAttachmentNotFound = errors.New("Attachment not found")

type AttachmentRepository interface {
	Create(arg entities.CreateAttachment) (entities.Attachment, error)
	Find(id string) (entities.Attachment, error)
	Delete(id string) error
}
