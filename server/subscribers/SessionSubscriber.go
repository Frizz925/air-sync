package subscribers

import (
	"air-sync/subscribers/events"
	"air-sync/util/pubsub"

	log "github.com/sirupsen/logrus"
)

type sessionSubscriber struct {
	stream *pubsub.Stream
}

func SubscribeSession(stream *pubsub.Stream) {
	sub := &sessionSubscriber{
		stream: stream,
	}
	sub.handleAsync(events.SessionDeleted, sub.handleSessionDelete)
	sub.handleAsync(events.MessageInserted, sub.handleMessageInsert)
	sub.handleAsync(events.MessageDeleted, sub.handleMessageDelete)
}

func (s *sessionSubscriber) handleAsync(name string, handler func(v interface{}) error) {
	go func() {
		err := s.stream.Topic(name).ForEach(handler)
		if err != nil {
			s.handleError(err)
		}
	}()
}

func (s *sessionSubscriber) handleSessionDelete(v interface{}) error {
	if evt, ok := v.(events.SessionDelete); ok {
		s.sessionTopic(string(evt)).Shutdown()
	}
	return nil
}

func (s *sessionSubscriber) handleMessageInsert(v interface{}) error {
	if evt, ok := v.(events.MessageInsert); ok {
		s.sessionTopic(evt.SessionId).Fire(events.SessionEvent{
			Event: events.MessageInserted,
			Value: evt.Message,
		})
	}
	return nil
}

func (s *sessionSubscriber) handleMessageDelete(v interface{}) error {
	if evt, ok := v.(events.MessageDelete); ok {
		s.sessionTopic(evt.SessionId).Fire(events.SessionEvent{
			Event: events.MessageDeleted,
			Value: evt.MessageId,
		})
	}
	return nil
}

func (s *sessionSubscriber) handleError(err error) {
	log.Error(err)
}

func (s *sessionSubscriber) sessionTopic(id string) *pubsub.Topic {
	return s.stream.Topic(events.SessionEventName(id))
}
