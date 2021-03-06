package services

import (
	"air-sync/models/events"
	"air-sync/util/pubsub"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

type EventService string

const (
	EventServiceRedis  EventService = "redis"
	EventServicePubSub EventService = "pubsub"
)

type RedisOptions struct {
	Addr     string
	Password string
}

type GooglePubSubOptions struct {
	ProjectID      string
	TopicID        string
	SubscriptionID string
}

type EventBrokerOptions struct {
	Service      EventService
	Redis        RedisOptions
	GooglePubSub GooglePubSubOptions
}

type EventBrokerService struct {
	EventBrokerOptions
	context     context.Context
	pub         *pubsub.Publisher
	broker      interface{}
	initialized bool
}

var _ Initializer = (*EventBrokerService)(nil)

func NewEventBrokerService(ctx context.Context, opts EventBrokerOptions) *EventBrokerService {
	return &EventBrokerService{
		EventBrokerOptions: opts,
		context:            ctx,
		pub:                pubsub.NewPublisher(),
		initialized:        false,
	}
}

func (s *EventBrokerService) Initialize() error {
	if s.initialized {
		return ErrAlreadyInitialized
	}
	s.pub.Topic(events.EventSession).Subscribe().
		ForEachAsync(s.context, s.handleSessionEvent, s.handleError)
	switch s.Service {
	case EventServiceRedis:
		s.broker = NewRedisBrokerService(s.context, RedisBrokerOptions{
			Publisher: s.pub,
			Addr:      s.Redis.Addr,
			Password:  s.Redis.Password,
		})
	case EventServicePubSub:
		s.broker = NewGooglePubSubBrokerService(s.context, GooglePubSubBrokerOptions{
			Publisher:      s.pub,
			ProjectID:      s.GooglePubSub.ProjectID,
			TopicID:        s.GooglePubSub.TopicID,
			SubscriptionID: s.GooglePubSub.SubscriptionID,
		})
	}
	if v, ok := s.broker.(Initializer); ok {
		if err := v.Initialize(); err != nil {
			return err
		}
	}
	s.initialized = true
	return nil
}

func (s *EventBrokerService) Deinitialize() {
	if !s.initialized {
		return
	}
	s.pub.Topic(events.EventSession).Close()
	if v, ok := s.broker.(Initializer); ok {
		v.Deinitialize()
	}
	s.initialized = false
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
	s.pub.Topic(events.EventSessionID(event.SessionID)).Publish(event)

	switch event.Event {
	case events.EventSessionDeleted:
		// Give grace period of 30 seconds before closing the topic
		go func() {
			time.Sleep(30 * time.Second)
			s.pub.Topic(events.EventSessionID(event.SessionID)).Close()
		}()
	}

	return nil
}

func (s *EventBrokerService) handleError(err error) {
	log.Error(err)
}
