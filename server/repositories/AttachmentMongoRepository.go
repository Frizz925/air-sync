package repositories

import (
	"air-sync/models"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotImplemented = errors.New("Not implemented")

type AttachmentMongoRepository struct {
	*MongoRepository
	attachments *mongo.Collection
}

var _ AttachmentRepository = (*AttachmentMongoRepository)(nil)

func NewAttachmentMongoRepository(db *mongo.Database) *AttachmentMongoRepository {
	return &AttachmentMongoRepository{
		MongoRepository: NewMongoRepository(db),
		attachments:     db.Collection("attachments"),
	}
}

func (r *AttachmentMongoRepository) Create(arg models.CreateAttachment) (models.Attachment, error) {
	return models.EmptyAttachment, ErrNotImplemented
}

func (r *AttachmentMongoRepository) Find(id string) (models.Attachment, error) {
	return models.EmptyAttachment, ErrNotImplemented
}

func (r *AttachmentMongoRepository) Delete(id string) error {
	return ErrNotImplemented
}
