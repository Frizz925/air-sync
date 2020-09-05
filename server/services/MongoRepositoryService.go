package services

import (
	repos "air-sync/repositories"
	"context"
	"net/url"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepositoryOptions struct {
	URL      *url.URL
	Database string
	Recreate bool
}

type MongoRepositoryService struct {
	context              context.Context
	client               *mongo.Client
	url                  *url.URL
	database             string
	sessionRepository    *repos.SessionMongoRepository
	attachmentRepository *repos.AttachmentMongoRepository
	recreate             bool
	initialized          bool
}

var _ RepositoryService = (*MongoRepositoryService)(nil)

func NewMongoRepositoryService(ctx context.Context, opts MongoRepositoryOptions) *MongoRepositoryService {
	return &MongoRepositoryService{
		context:     context.Background(),
		url:         opts.URL,
		database:    opts.Database,
		recreate:    opts.Recreate,
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
	log.Infof("Using MongoDB database: %s", s.database)
	db := client.Database(s.database)
	opts := repos.MongoOptions{
		Database: db,
		Recreate: s.recreate,
	}

	sessionRepo := repos.NewSessionMongoRepository(s.context, opts)
	if err := sessionRepo.Migrate(); err != nil {
		return err
	}
	s.sessionRepository = sessionRepo

	attachmentRepo := repos.NewAttachmentMongoRepository(s.context, opts)
	if err := attachmentRepo.Migrate(); err != nil {
		return err
	}
	s.attachmentRepository = attachmentRepo

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
