package repositories

import (
	"air-sync/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotImplemented = errors.New("Not implemented")

type AttachmentMongoRepository struct {
	*MongoRepository
	context     context.Context
	attachments *mongo.Collection
}

var _ AttachmentRepository = (*AttachmentMongoRepository)(nil)

func NewAttachmentMongoRepository(ctx context.Context, db *mongo.Database) *AttachmentMongoRepository {
	return &AttachmentMongoRepository{
		MongoRepository: NewMongoRepository(db),
		context:         ctx,
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
