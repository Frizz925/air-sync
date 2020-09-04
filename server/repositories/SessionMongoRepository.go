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
	session := mongoModels.Session{}
	{
		cur, err := r.sessions.Find(r.context, bson.M{"id": id})
		if err != nil {
			return models.EmptySession, err
		}
		defer cur.Close(r.context)
		if !cur.Next(r.context) {
			return models.EmptySession, ErrSessionNotFound
		}
		if err := cur.Decode(&session); err != nil {
			return models.EmptySession, err
		}
	}
	messages := make([]mongoModels.Message, 0)
	{
		cur, err := r.messages.Find(r.context, bson.M{"session_id": id})
		if err != nil {
			return models.EmptySession, err
		}
		defer cur.Close(r.context)
		if err := cur.All(r.context, &messages); err != nil {
			return models.EmptySession, err
		}
	}
	messageResults := make([]models.Message, len(messages))
	{
		attachmentSet := make(map[string]bool)
		attachmentIds := make(bson.A, 0)
		for _, message := range messages {
			attachmentId := message.AttachmentID
			if attachmentId == "" {
				continue
			}
			attachmentSet[attachmentId] = true
			if _, ok := attachmentSet[attachmentId]; !ok {
				attachmentIds = append(attachmentIds, attachmentId)
			}
		}
		attachmentMap, err := r.FindAttachments(attachmentIds)
		if err != nil {
			return models.EmptySession, err
		}
		for idx, message := range messages {
			attachment := attachmentMap[message.AttachmentID]
			messageResults[idx] = mongoModels.ToMessageModel(message, attachment)
		}
	}
	return mongoModels.ToSessionModel(session, messageResults), nil
}

func (r *SessionMongoRepository) FindBefore(t time.Time) ([]models.Session, error) {
	sessions := make([]models.Session, 0)
	cur, err := r.sessions.Find(r.context, bson.M{
		"created_at": bson.M{"$lt": t.Unix()},
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
	message := mongoModels.FromInsertMessageModel(id, arg)
	if _, err := r.messages.InsertOne(r.context, message); err != nil {
		return models.EmptyMessage, err
	}
	if message.AttachmentID == "" {
		return mongoModels.ToMessageModel(message, mongoModels.EmptyAttachment), nil
	}
	attachment, err := r.FindOneAttachment(message.AttachmentID)
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
	attachmentMap, err := r.FindAttachments(bson.A{id})
	return attachmentMap[id], err
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
