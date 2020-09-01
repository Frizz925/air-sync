package services

import (
	repos "air-sync/repositories"
	"context"
	"net/url"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepositoryService struct {
	context              context.Context
	client               *mongo.Client
	url                  *url.URL
	dbName               string
	sessionRepository    *repos.SessionMongoRepository
	attachmentRepository *repos.AttachmentMongoRepository
	initialized          bool
}

var _ RepositoryService = (*MongoRepositoryService)(nil)

func NewMongoRepositoryService(url *url.URL, dbName string) *MongoRepositoryService {
	return &MongoRepositoryService{
		context:     context.Background(),
		url:         url,
		dbName:      dbName,
		initialized: false,
	}
}

func (s *MongoRepositoryService) Initialize() error {
	if s.initialized {
		return ErrAlreadyInitialized
	}

	log.Infof("Connecting to MongoDB: %s", s.url.Host)
	client, err := mongo.Connect(s.context, options.Client().ApplyURI(s.url.String()))
	if err != nil {
		return err
	}
	s.client = client
	if err := client.Ping(s.context, nil); err != nil {
		defer s.disconnect()
		return err
	}
	log.Infof("Connected to MongoDB")

	log.Infof("Using MongoDB database: %s", s.dbName)
	db := client.Database(s.dbName)
	s.sessionRepository = repos.NewSessionMongoRepository(s.context, db)
	s.attachmentRepository = repos.NewAttachmentMongoRepository(s.context, db)

	s.initialized = true
	return nil
}

func (s *MongoRepositoryService) Deinitialize() {
	if !s.initialized {
		log.Error(ErrNotInitialized)
		return
	}
	if err := s.client.Disconnect(s.context); err != nil {
		log.Error(err)
	}
}

func (s *MongoRepositoryService) SessionRepository() repos.SessionRepository {
	return s.sessionRepository
}

func (s *MongoRepositoryService) AttachmentRepository() repos.AttachmentRepository {
	return s.attachmentRepository
}

func (s *MongoRepositoryService) disconnect() {
	if s.client != nil {
		err := s.client.Disconnect(s.context)
		if err != nil {
			log.Error(err)
		}
	}
}
