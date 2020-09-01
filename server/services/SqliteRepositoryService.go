package services

import (
	repos "air-sync/repositories"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteRepositoryService struct {
	dsn               string
	sessionRepository *repos.SessionSqlRepository
}

var _ RepositoryService = (*SqliteRepositoryService)(nil)

func NewSqliteRepositoryService(dsn string) *SqliteRepositoryService {
	return &SqliteRepositoryService{
		dsn: dsn,
	}
}

func (s *SqliteRepositoryService) Initialize() error {
	db, err := gorm.Open(sqlite.Open(s.dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sessionRepo := repos.NewSessionSqlRepository(db)
	if err := sessionRepo.Migrate(); err != nil {
		return err
	}
	s.sessionRepository = sessionRepo

	return nil
}

func (s *SqliteRepositoryService) SessionRepository() repos.SessionRepository {
	return s.sessionRepository
}
