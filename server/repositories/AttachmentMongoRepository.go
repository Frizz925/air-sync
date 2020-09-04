package repositories

import (
	"air-sync/models"
	mongoModels "air-sync/models/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const MongoAttachmentCollection = "attachments"

type AttachmentMongoRepository struct {
	*MongoRepository
	context     context.Context
	attachments *mongo.Collection
}

var _ AttachmentRepository = (*AttachmentMongoRepository)(nil)
var _ RepositoryMigration = (*AttachmentMongoRepository)(nil)

func NewAttachmentMongoRepository(ctx context.Context, db *mongo.Database) *AttachmentMongoRepository {
	return &AttachmentMongoRepository{
		MongoRepository: NewMongoRepository(db),
		context:         ctx,
		attachments:     db.Collection(MongoAttachmentCollection),
	}
}

func (r *AttachmentMongoRepository) Migrate() error {
	_, err := r.attachments.Indexes().CreateMany(r.context, []mongo.IndexModel{
		{Keys: bson.M{"id": "hashed"}},
		{Keys: bson.M{"created_at": 1}},
	})
	return err
}

func (r *AttachmentMongoRepository) Create(arg models.CreateAttachment) (models.Attachment, error) {
	attachment := mongoModels.FromCreateAttachmentModel(arg)
	_, err := r.attachments.InsertOne(r.context, attachment)
	return mongoModels.ToAttachmentModel(attachment), err
}

func (r *AttachmentMongoRepository) Find(id string) (models.Attachment, error) {
	cur, err := r.attachments.Find(r.context, bson.M{"id": id})
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

func (r *AttachmentMongoRepository) FindOrphans() ([]models.Attachment, error) {
	attachments := make([]models.Attachment, 0)
	cur, err := r.attachments.Aggregate(r.context, bson.A{
		bson.M{"$lookup": bson.M{
			"from":         "messages",
			"localField":   "id",
			"foreignField": "attachment_id",
			"as":           "messages",
		}},
		bson.M{"$match": bson.M{"messages": bson.A{}}},
	})
	if err != nil {
		return attachments, err
	}
	defer cur.Close(r.context)
	for cur.Next(r.context) {
		attachment := mongoModels.Attachment{}
		if err := cur.Decode(&attachment); err != nil {
			return attachments, err
		}
		attachments = append(attachments, mongoModels.ToAttachmentModel(attachment))
	}
	return attachments, nil
}

func (r *AttachmentMongoRepository) Delete(id string) error {
	res, err := r.attachments.DeleteOne(r.context, bson.M{"id": id})
	if err != nil {
		return err
	} else if res.DeletedCount <= 0 {
		return ErrAttachmentNotFound
	}
	return nil
}

func (r *AttachmentMongoRepository) DeleteMany(ids []string) (int, error) {
	res, err := r.attachments.DeleteMany(r.context, bson.M{
		"id": bson.M{"$in": ids},
	})
	return int(res.DeletedCount), err
}
