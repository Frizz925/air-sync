package repositories

import (
	"air-sync/models"
	mongoModels "air-sync/models/mongo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MongoSessionCollection = "sessions"
	MongoMessageCollection = "messages"
)

type SessionMongoRepository struct {
	*MongoRepository
	context     context.Context
	sessions    *mongo.Collection
	messages    *mongo.Collection
	attachments *mongo.Collection
}

type mongoMessageQuery struct {
	mongoModels.Message `bson:"inline"`
	Attachment          mongoModels.Attachment `bson:"attachment"`
}

type mongoSessionQuery struct {
	mongoModels.Session `bson:"inline"`
	Messages            []mongoMessageQuery `bson:"messages"`
}

var _ SessionRepository = (*SessionMongoRepository)(nil)
var _ RepositoryMigration = (*SessionMongoRepository)(nil)

func NewSessionMongoRepository(ctx context.Context, db *mongo.Database) *SessionMongoRepository {
	return &SessionMongoRepository{
		MongoRepository: NewMongoRepository(db),
		context:         ctx,
		sessions:        db.Collection(MongoSessionCollection),
		messages:        db.Collection(MongoMessageCollection),
		attachments:     db.Collection(MongoAttachmentCollection),
	}
}

func (r *SessionMongoRepository) Migrate() error {
	{
		_, err := r.sessions.Indexes().CreateMany(r.context, []mongo.IndexModel{
			{Keys: bson.M{"id": "hashed"}},
			{Keys: bson.M{"created_at": 1}},
		})
		if err != nil {
			return err
		}
	}
	{
		_, err := r.messages.Indexes().CreateMany(r.context, []mongo.IndexModel{
			{Keys: bson.M{"id": "hashed"}},
			{Keys: bson.M{"session_id": "hashed"}},
			{Keys: bson.M{"attachment_id": "hashed"}},
			{Keys: bson.M{"created_at": 1}},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SessionMongoRepository) Create() (models.Session, error) {
	session := mongoModels.NewSession()
	_, err := r.sessions.InsertOne(r.context, session)
	messages := make([]models.Message, 0)
	return mongoModels.ToSessionModel(session, messages), err
}

func (r *SessionMongoRepository) Find(id string) (models.Session, error) {
	cur, err := r.sessions.Aggregate(r.context, bson.A{
		bson.M{"$match": bson.M{"id": id}},
		bson.M{"$lookup": bson.M{
			"from": MongoMessageCollection,
			"let":  bson.M{"session_id": "$id"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{
					"$eq": bson.A{"$session_id", "$$session_id"},
				}}},
				bson.M{"$sort": bson.M{"created_at": -1}},
				bson.M{"$lookup": bson.M{
					"from":         MongoAttachmentCollection,
					"localField":   "attachment_id",
					"foreignField": "id",
					"as":           "attachment",
				}},
				bson.M{"$unwind": bson.M{
					"path":                       "$attachment",
					"preserveNullAndEmptyArrays": true,
				}},
			},
			"as": "messages",
		}},
	})
	if err != nil {
		return models.EmptySession, err
	}
	defer cur.Close(r.context)
	if !cur.Next(r.context) {
		return models.EmptySession, ErrSessionNotFound
	}
	session := mongoSessionQuery{}
	if err := cur.Decode(&session); err != nil {
		return models.EmptySession, err
	}
	messages := make([]models.Message, len(session.Messages))
	for idx, message := range session.Messages {
		messages[idx] = mongoModels.ToMessageModel(message.Message, message.Attachment)
	}
	return mongoModels.ToSessionModel(session.Session, messages), nil
}

func (r *SessionMongoRepository) FindBefore(t time.Time) ([]models.Session, error) {
	sessions := make([]models.Session, 0)
	cur, err := r.sessions.Find(r.context, bson.M{
		"created_at": bson.M{"$lt": models.FromTime(t)},
	})
	if err != nil {
		return sessions, err
	}
	defer cur.Close(r.context)
	for cur.Next(r.context) {
		session := mongoModels.Session{}
		if err := cur.Decode(&session); err != nil {
			return sessions, err
		}
		sessions = append(sessions, mongoModels.ToSessionModel(session, nil))
	}
	return sessions, nil
}

func (r *SessionMongoRepository) InsertMessage(id string, arg models.InsertMessage) (models.Message, error) {
	cur, err := r.sessions.Find(r.context, bson.M{"id": id})
	if err != nil {
		return models.EmptyMessage, err
	}
	defer cur.Close(r.context)
	if !cur.TryNext(r.context) {
		return models.EmptyMessage, ErrSessionNotFound
	}
	attachment := mongoModels.EmptyAttachment
	if arg.AttachmentID != "" {
		res, err := r.FindOneAttachment(arg.AttachmentID)
		if err != nil {
			return models.EmptyMessage, err
		}
		attachment = res
	}
	message := mongoModels.FromInsertMessageModel(id, arg)
	if _, err := r.messages.InsertOne(r.context, message); err != nil {
		return models.EmptyMessage, err
	}
	return mongoModels.ToMessageModel(message, attachment), err
}

func (r *SessionMongoRepository) DeleteMessage(id string, messageID string) error {
	res, err := r.messages.DeleteOne(r.context, bson.M{"id": messageID, "session_id": id})
	if err != nil {
		return err
	} else if res.DeletedCount <= 0 {
		return ErrMessageNotFound
	}
	return nil
}

func (r *SessionMongoRepository) FindOneAttachment(id string) (mongoModels.Attachment, error) {
	cur, err := r.attachments.Find(r.context, bson.M{"id": id})
	if err != nil {
		return mongoModels.EmptyAttachment, err
	}
	defer cur.Close(r.context)
	if !cur.Next(r.context) {
		return mongoModels.EmptyAttachment, ErrAttachmentNotFound
	}
	attachment := mongoModels.Attachment{}
	return attachment, cur.Decode(&attachment)
}

func (r *SessionMongoRepository) FindAttachments(ids bson.A) (map[string]mongoModels.Attachment, error) {
	resultMap := make(map[string]mongoModels.Attachment)
	cur, err := r.attachments.Find(r.context, bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return resultMap, err
	}
	defer cur.Close(r.context)
	for cur.Next(r.context) {
		attachment := mongoModels.Attachment{}
		if err := cur.Decode(&attachment); err != nil {
			return resultMap, err
		}
		resultMap[attachment.ID] = attachment
	}
	return resultMap, nil
}

func (r *SessionMongoRepository) Delete(id string) error {
	_, err := r.messages.DeleteMany(r.context, bson.M{"session_id": id})
	if err != nil {
		return err
	}
	res, err := r.sessions.DeleteOne(r.context, bson.M{"id": id})
	if err != nil {
		return err
	} else if res.DeletedCount <= 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (r *SessionMongoRepository) DeleteMany(ids []string) (int, error) {
	_, err := r.messages.DeleteMany(r.context, bson.M{
		"session_id": bson.M{"$in": ids},
	})
	if err != nil {
		return 0, err
	}
	res, err := r.sessions.DeleteMany(r.context, bson.M{
		"id": bson.M{"$in": ids},
	})
	return int(res.DeletedCount), err
}
