package services

import (
	repos "air-sync/repositories"
)

type RepositoryService interface {
	SessionRepository() repos.SessionRepository
	AttachmentRepository() repos.AttachmentRepository
}
