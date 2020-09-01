package services

import (
	repos "air-sync/repositories"

	"gorm.io/gorm"
)

type GormRepositoryService struct {
	dialector            gorm.Dialector
	sessionRepository    *repos.SessionSqlRepository
	attachmentRepository *repos.AttachmentSqlRepository
	initialized          bool
}

var _ RepositoryService = (*GormRepositoryService)(nil)

func NewGormRepositoryService(dialector gorm.Dialector) *GormRepositoryService {
	return &GormRepositoryService{
		dialector:   dialector,
		initialized: false,
	}
}

func (s *GormRepositoryService) Initialize() error {
	if s.initialized {
		return ErrAlreadyInitialized
	}

	db, err := gorm.Open(s.dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	sessionRepo := repos.NewSessionSqlRepository(db)
	if err := sessionRepo.Migrate(); err != nil {
		return err
	}
	s.sessionRepository = sessionRepo

	attachmentRepo := repos.NewAttachmentSqlRepository(db)
	if err := attachmentRepo.Migrate(); err != nil {
		return err
	}
	s.attachmentRepository = attachmentRepo

	s.initialized = true
	return nil
}

func (s *GormRepositoryService) SessionRepository() repos.SessionRepository {
	return s.sessionRepository
}

func (s *GormRepositoryService) AttachmentRepository() repos.AttachmentRepository {
	return s.attachmentRepository
}
