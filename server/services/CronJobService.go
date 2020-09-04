package services

import (
	repos "air-sync/repositories"
	"air-sync/storages"
	"air-sync/util/pubsub"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type CronJobOptions struct {
	SessionRepository    repos.SessionRepository
	AttachmentRepository repos.AttachmentRepository
	Publisher            *pubsub.Publisher
	Storage              storages.Storage
}

type CronJobService struct {
	sessionRepo    repos.SessionRepository
	attachmentRepo repos.AttachmentRepository
	pub            *pubsub.Publisher
	storage        storages.Storage
}

func NewCronJobService(opts CronJobOptions) *CronJobService {
	return &CronJobService{
		sessionRepo:    opts.SessionRepository,
		attachmentRepo: opts.AttachmentRepository,
		pub:            opts.Publisher,
		storage:        opts.Storage,
	}
}

func (s *CronJobService) RunCleanupJob() error {
	{
		s.log("Deleting old sessions")
		n, err := s.sessionRepo.DeleteBefore(time.Now().Add(-24 * time.Hour))
		if err != nil {
			return err
		}
		s.log("Deleted %d session(s)", n)
	}
	{
		s.log("Deleting orphan attachments")
		attachments, err := s.attachmentRepo.FindOrphans()
		if err != nil {
			return err
		}
		attachmentIds := make([]string, len(attachments))
		for _, attachment := range attachments {
			attachmentIds = append(attachmentIds, attachment.ID)
		}
		n, err := s.attachmentRepo.DeleteMany(attachmentIds)
		if err != nil {
			return err
		}
		s.log("Deleted %d attachment(s)", n)
	}
	return nil
}

func (s *CronJobService) log(format string, a ...interface{}) {
	log.Info("Cron: " + fmt.Sprintf(format, a...))
}
