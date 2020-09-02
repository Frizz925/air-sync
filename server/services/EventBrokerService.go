package services

import (
	"air-sync/models/events"
	"air-sync/util/pubsub"
	"context"

	log "github.com/sirupsen/logrus"
)

type EventBrokerService struct {
	context context.Context
	pub     *pubsub.Publisher
}

func NewEventBrokerService(ctx context.Context) *EventBrokerService {
	return &EventBrokerService{
		context: ctx,
		pub:     pubsub.NewPublisher(),
	}
}

func (s *EventBrokerService) Initialize() {
	s.pub.Topic(events.EventSession).Subscribe().
		ForEachAsync(s.context, s.handleSessionEvent, s.handleError)
}

func (s *EventBrokerService) Deinitialize() {
	s.pub.Topic(events.EventSession).Close()
}

func (s *EventBrokerService) Publisher() *pubsub.Publisher {
	return s.pub
}

func (s *EventBrokerService) handleSessionEvent(v interface{}) error {
	event, ok := v.(events.SessionEvent)
	if !ok {
		return nil
	}
	if event.Error != nil {
		return event.Error
	}

	s.pub.Topic(event.Event).Publish(event.Value)
	s.pub.Topic(events.EventSessionId(event.SessionId)).Publish(event)

	// post-publish
	switch event.Event {
	case events.EventSessionDeleted:
		s.pub.Topic(events.EventSessionId(event.SessionId)).Close()
	}

	return nil
}

func (s *EventBrokerService) handleError(err error) {
	log.Error(err)
}
