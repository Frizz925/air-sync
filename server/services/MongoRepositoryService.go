package services

import (
	repos "air-sync/repositories"
	"context"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepositoryService struct {
	client               *mongo.Client
	mongoUri             string
	dbName               string
	sessionRepository    *repos.SessionMongoRepository
	attachmentRepository *repos.AttachmentMongoRepository
}

var _ RepositoryService = (*MongoRepositoryService)(nil)

func NewMongoRepositoryService(mongoUri string, dbName string) *MongoRepositoryService {
	return &MongoRepositoryService{
		mongoUri: mongoUri,
		dbName:   dbName,
	}
}

func (s *MongoRepositoryService) Initialize() error {
	log.Infof("Connecting to MongoDB: %s", s.mongoUri)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(s.mongoUri))
	if err != nil {
		return err
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}
	log.Infof("Connected to MongoDB")
	s.client = client

	log.Infof("Using MongoDB database: %s", s.dbName)
	db := client.Database(s.dbName)
	s.sessionRepository = repos.NewSessionMongoRepository(db)
	s.attachmentRepository = repos.NewAttachmentMongoRepository(db)

	return nil
}

func (s *MongoRepositoryService) Deinitialize() {
	if err := s.client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (s *MongoRepositoryService) SessionRepository() repos.SessionRepository {
	return s.sessionRepository
}

func (s *MongoRepositoryService) AttachmentRepository() repos.AttachmentRepository {
	return s.attachmentRepository
}
