package repositories

import (
	"air-sync/models"
	mongoModels "air-sync/models/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionMongoRepository struct {
	*MongoRepository
	sessions *mongo.Collection
}

var _ SessionRepository = (*SessionMongoRepository)(nil)

func NewSessionMongoRepository(db *mongo.Database) *SessionMongoRepository {
	return &SessionMongoRepository{
		MongoRepository: NewMongoRepository(db),
		sessions:        db.Collection("sessions"),
	}
}

func (r *SessionMongoRepository) Create() (models.Session, error) {
	session := mongoModels.NewSession()
	_, err := r.sessions.InsertOne(context.Background(), session)
	return mongoModels.ToSessionModel(session), err
}

func (r *SessionMongoRepository) Find(id string) (models.Session, error) {
	cur, err := r.sessions.Find(context.Background(), bson.M{"_id": id})
	if err != nil {
		return models.EmptySession, err
	}
	defer cur.Close(context.Background())
	if !cur.Next(context.Background()) {
		return models.EmptySession, ErrSessionNotFound
	}
	session := mongoModels.Session{}
	err = cur.Decode(&session)
	return mongoModels.ToSessionModel(session), err
}

func (r *SessionMongoRepository) InsertMessage(id string, arg models.InsertMessage) (models.Message, error) {
	message := mongoModels.FromInsertMessageModel(arg)
	res, err := r.sessions.UpdateOne(
		context.Background(),
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
	return mongoModels.ToMessageModel(message), nil
}

func (r *SessionMongoRepository) DeleteMessage(id string, messageId string) error {
	cur, err := r.sessions.Find(
		context.Background(),
		bson.M{"_id": id, "messages._id": messageId},
	)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	if !cur.Next(context.Background()) {
		return ErrMessageNotFound
	}
	_, err = r.sessions.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"messages": bson.M{"_id": messageId}}},
	)
	return err
}

func (r *SessionMongoRepository) Delete(id string) error {
	res, err := r.sessions.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	} else if res.DeletedCount <= 0 {
		return ErrSessionNotFound
	}
	return nil
}
