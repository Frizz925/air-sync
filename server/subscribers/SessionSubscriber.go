package subscribers

import (
	repos "air-sync/repositories"
	"air-sync/subscribers/events"
	"air-sync/util/pubsub"

	log "github.com/sirupsen/logrus"
)

type sessionSubscriber struct {
	repo   repos.SessionRepository
	stream *pubsub.Stream
}

func SubscribeSession(repo repos.SessionRepository, stream *pubsub.Stream) {
	sub := &sessionSubscriber{
		repo:   repo,
		stream: stream,
	}
	sub.handleAsync(events.SessionUpdateEventName, sub.handleUpdate)
	sub.handleAsync(events.SessionDeleteEventName, sub.handleDelete)
}

func (s *sessionSubscriber) handleAsync(name string, handler func(v interface{}) error) {
	go func() {
		err := s.stream.Topic(name).ForEach(handler)
		if err != nil {
			s.handleError(err)
		}
	}()
}

func (s *sessionSubscriber) handleUpdate(v interface{}) error {
	if evt, ok := v.(events.SessionUpdate); ok {
		s.sessionTopic(evt.Id).Fire(evt.Message)
		return s.repo.Update(evt.Id, evt.Message)
	}
	return nil
}

func (s *sessionSubscriber) handleDelete(v interface{}) error {
	if evt, ok := v.(events.SessionDelete); ok {
		id := string(evt)
		s.sessionTopic(id).Shutdown()
		return s.repo.Delete(id)
	}
	return nil
}

func (s *sessionSubscriber) handleError(err error) {
	log.Error(err)
}

func (s *sessionSubscriber) sessionTopic(id string) *pubsub.Topic {
	return s.stream.Topic(events.SessionEventName(id))
}
