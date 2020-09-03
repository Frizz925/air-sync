package repositories

import (
	"air-sync/models"
	mongoModels "air-sync/models/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoAttachmentCollection = "attachments"

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
		attachments:     db.Collection(MongoAttachmentCollection),
	}
}

func (r *AttachmentMongoRepository) Create(arg models.CreateAttachment) (models.Attachment, error) {
	attachment := mongoModels.FromCreateAttachmentModel(arg)
	_, err := r.attachments.InsertOne(r.context, attachment)
	return mongoModels.ToAttachmentModel(attachment), err
}

func (r *AttachmentMongoRepository) Find(id string) (models.Attachment, error) {
	cur, err := r.attachments.Find(r.context, bson.M{"_id": id})
	if err != nil {
		return models.EmptyAttachment, err
	}
	defer cur.Close(r.context)
	if !cur.Next(r.context) {
		return models.EmptyAttachment, ErrAttachmentNotFound
	}
	attachment := mongoModels.Attachment{}
	if err := cur.Decode(&attachment); err != nil {
		return models.EmptyAttachment, err
	}
	return mongoModels.ToAttachmentModel(attachment), nil
}

func (r *AttachmentMongoRepository) Delete(id string) error {
	res, err := r.attachments.DeleteOne(r.context, bson.M{"_id": id})
	if err != nil {
		return err
	} else if res.DeletedCount <= 0 {
		return ErrAttachmentNotFound
	}
	return nil
}
