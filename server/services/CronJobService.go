package services

import (
	"air-sync/models/events"
	repos "air-sync/repositories"
	"air-sync/storages"
	"air-sync/util/pubsub"
	"fmt"
	"sync"
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
	topic          *pubsub.Topic
	storage        storages.Storage
	nextRun        time.Time
	interval       time.Duration
	mu             sync.Mutex
}

func NewCronJobService(opts CronJobOptions) *CronJobService {
	return &CronJobService{
		sessionRepo:    opts.SessionRepository,
		attachmentRepo: opts.AttachmentRepository,
		topic:          opts.Publisher.Topic(events.EventSession),
		storage:        opts.Storage,
		nextRun:        time.Unix(0, 0),
		interval:       1 * time.Hour,
	}
}

func (s *CronJobService) RunCleanupJob() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if time.Now().Before(s.nextRun) {
		dt := s.nextRun.UTC().Format(time.RFC3339)
		return NewCronRequestError("No cleanup job run until %s", dt)
	}
	deadline := time.Now().Add(-24 * time.Hour)
	{
		s.log("Deleting old sessions")
		sessions, err := s.sessionRepo.FindBefore(deadline)
		if err != nil {
			return err
		}
		sessionIds := make([]string, len(sessions))
		for idx, session := range sessions {
			sessionIds[idx] = session.ID
		}
		n, err := s.sessionRepo.DeleteMany(sessionIds)
		if err != nil {
			return err
		}
		for _, id := range sessionIds {
			s.topic.Publish(events.CreateSessionEvent(
				id, events.EventSessionDeleted,
				events.SessionDelete(id), nil,
			))
		}
		s.log("Deleted %d session(s)", n)
	}
	{
		s.log("Deleting orphan attachments")
		attachments, err := s.attachmentRepo.FindOrphansBefore(deadline)
		if err != nil {
			return err
		}
		attachmentIds := make([]string, len(attachments))
		for idx, attachment := range attachments {
			attachmentIds[idx] = attachment.ID
		}
		for _, id := range attachmentIds {
			if id == "" {
				continue
			}
			exists, err := s.storage.Exists(id)
			if err != nil {
				return err
			} else if !exists {
				continue
			}
			if err := s.storage.Delete(id); err != nil {
				return err
			}
		}
		n, err := s.attachmentRepo.DeleteMany(attachmentIds)
		if err != nil {
			return err
		}
		s.log("Deleted %d attachment(s)", n)
	}
	s.nextRun = time.Now().Add(s.interval)
	return nil
}

func (s *CronJobService) log(format string, a ...interface{}) {
	log.Info("Cron: " + fmt.Sprintf(format, a...))
}
