package repositories

import (
	"air-sync/models"
	mongoModels "air-sync/models/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoSessionCollection = "sessions"

type SessionMongoRepository struct {
	*MongoRepository
	context     context.Context
	sessions    *mongo.Collection
	attachments *mongo.Collection
}

var _ SessionRepository = (*SessionMongoRepository)(nil)

func NewSessionMongoRepository(ctx context.Context, db *mongo.Database) *SessionMongoRepository {
	return &SessionMongoRepository{
		MongoRepository: NewMongoRepository(db),
		context:         ctx,
		sessions:        db.Collection(MongoSessionCollection),
		attachments:     db.Collection(MongoAttachmentCollection),
	}
}

func (r *SessionMongoRepository) Create() (models.Session, error) {
	session := mongoModels.NewSession()
	_, err := r.sessions.InsertOne(r.context, session)
	messages := make([]models.Message, 0)
	return mongoModels.ToSessionModel(session, messages), err
}

func (r *SessionMongoRepository) Find(id string) (models.Session, error) {
	cur, err := r.sessions.Find(r.context, bson.M{"_id": id})
	if err != nil {
		return models.EmptySession, err
	}
	defer cur.Close(r.context)
	if !cur.Next(r.context) {
		return models.EmptySession, ErrSessionNotFound
	}
	session := mongoModels.Session{}
	if err := cur.Decode(&session); err != nil {
		return models.EmptySession, err
	}
	messages := make([]models.Message, len(session.Messages))
	if len(session.Messages) <= 0 {
		return mongoModels.ToSessionModel(session, messages), nil
	}
	attachmentSet := make(map[string]bool)
	for _, message := range session.Messages {
		if message.AttachmentID != "" {
			attachmentSet[message.AttachmentID] = true
		}
	}
	attachmentIds := make(bson.A, 0)
	for attachmentId := range attachmentSet {
		attachmentIds = append(attachmentIds, attachmentId)
	}
	attachmentMap, err := r.FindAttachments(attachmentIds)
	if err != nil {
		return models.EmptySession, err
	}
	for idx, message := range session.Messages {
		attachment := attachmentMap[message.AttachmentID]
		messages[idx] = mongoModels.ToMessageModel(message, attachment)
	}
	return mongoModels.ToSessionModel(session, messages), nil
}

func (r *SessionMongoRepository) InsertMessage(id string, arg models.InsertMessage) (models.Message, error) {
	message := mongoModels.FromInsertMessageModel(arg)
	res, err := r.sessions.UpdateOne(
		r.context,
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"messages": bson.M{
			"$each":     bson.A{message},
			"$position": 0,
		}}},
	)
	if err != nil {
		return models.EmptyMessage, err
	} else if res.MatchedCount <= 0 {
		return models.EmptyMessage, ErrSessionNotFound
	}
	if message.AttachmentID == "" {
		return mongoModels.ToMessageModel(message, mongoModels.EmptyAttachment), nil
	}
	attachment, err := r.FindOneAttachment(message.AttachmentID)
	return mongoModels.ToMessageModel(message, attachment), err
}

func (r *SessionMongoRepository) DeleteMessage(id string, messageID string) error {
	cur, err := r.sessions.Find(
		r.context,
		bson.M{"_id": id, "messages._id": messageID},
	)
	if err != nil {
		return err
	}
	defer cur.Close(r.context)
	if !cur.Next(r.context) {
		return ErrMessageNotFound
	}
	_, err = r.sessions.UpdateOne(
		r.context,
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"messages": bson.M{"_id": messageID}}},
	)
	return err
}

func (r *SessionMongoRepository) FindOneAttachment(id string) (mongoModels.Attachment, error) {
	attachmentMap, err := r.FindAttachments(bson.A{id})
	return attachmentMap[id], err
}

func (r *SessionMongoRepository) FindAttachments(ids bson.A) (map[string]mongoModels.Attachment, error) {
	resultMap := make(map[string]mongoModels.Attachment)
	cur, err := r.attachments.Find(r.context, bson.M{"_id": bson.M{"$in": ids}})
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
	res, err := r.sessions.DeleteOne(r.context, bson.M{"_id": id})
	if err != nil {
		return err
	} else if res.DeletedCount <= 0 {
		return ErrSessionNotFound
	}
	return nil
}
