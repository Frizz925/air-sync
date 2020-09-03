package services

import (
	"air-sync/models/events"
	"air-sync/util/pubsub"
	"context"
	"encoding/json"
	"sync"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

type RedisBrokerOptions struct {
	Publisher *pubsub.Publisher
	Addr      string
	Password  string
}

type RedisBrokerService struct {
	context context.Context
	mu      sync.Mutex

	pub *pubsub.Publisher

	client   *redis.Client
	addr     string
	password string

	clientID    string
	lastEventID string

	initialized bool
}

var _ Service = (*RedisBrokerService)(nil)

func NewRedisBrokerService(ctx context.Context, opts RedisBrokerOptions) *RedisBrokerService {
	return &RedisBrokerService{
		context:     ctx,
		pub:         opts.Publisher,
		addr:        opts.Addr,
		password:    opts.Password,
		clientID:    uuid.NewV4().String(),
		initialized: false,
	}
}

func (s *RedisBrokerService) Initialize() error {
	if s.initialized {
		return ErrAlreadyInitialized
	}

	log.Infof("Connecting to Redis: %s", s.addr)
	client := redis.NewClient(&redis.Options{
		Addr:     s.addr,
		Password: s.password,
	})
	if err := client.Ping(s.context).Err(); err != nil {
		return err
	}
	s.client = client
	log.Info("Connected to Redis")

	s.handlePublishingAsync()
	s.handleSubscriptionAsync()
	log.Infof("Publishing and subscribing to Redis PubSub: %s", events.EventSession)

	s.initialized = true
	return nil
}

func (s *RedisBrokerService) Deinitialize() {
	if !s.initialized {
		log.Error(ErrNotInitialized)
		return
	}
	s.initialized = false
}

func (s *RedisBrokerService) handlePublishingAsync() {
	s.pub.Topic(events.EventSession).Subscribe().
		ForEachAsync(s.context, s.handlePublishing, s.handleError)
}

func (s *RedisBrokerService) handlePublishing(v interface{}) error {
	event, ok := v.(events.SessionEvent)
	if !ok {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID == s.lastEventID {
		return nil
	}
	s.lastEventID = event.ID
	b, err := json.Marshal(events.PubSubSessionEvent{
		SessionEvent: event,
		ClientID:     s.clientID,
	})
	if err != nil {
		return err
	}
	res := s.client.Publish(s.context, events.EventSession, string(b))
	if err := res.Err(); err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"id":        event.ID,
		"event":     event.Event,
		"timestamp": event.Timestamp,
	}).Infof("Redis published event")
	return nil
}

func (s *RedisBrokerService) handleSubscriptionAsync() {
	go func() {
		err := s.handleSubscription()
		if err != nil {
			s.handleError(err)
		}
	}()
}

func (s *RedisBrokerService) handleSubscription() error {
	ps := s.client.Subscribe(s.context, events.EventSession)
	if err := ps.Ping(s.context); err != nil {
		return err
	}
	if _, err := ps.Receive(s.context); err != nil {
		return err
	}
	ch := ps.Channel()
	for {
		select {
		case msg := <-ch:
			err := s.handleSubscriptionMessage(msg)
			if err != nil {
				return err
			}
		case <-s.context.Done():
			return nil
		}
	}
}

func (s *RedisBrokerService) handleSubscriptionMessage(msg *redis.Message) error {
	if msg.Channel != events.EventSession {
		return nil
	}
	payload := []byte(msg.Payload)
	event := events.PubSubSessionEvent{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	// Prevent pubsub self-loop
	if event.ClientID == s.clientID {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID == s.lastEventID {
		return nil
	}
	s.lastEventID = event.ID
	log.WithFields(log.Fields{
		"id":        event.ID,
		"event":     event.Event,
		"timestamp": event.Timestamp,
	}).Infof("Redis received event")
	s.pub.Topic(events.EventSession).Publish(event.SessionEvent)
	return nil
}

func (s *RedisBrokerService) handleError(err error) {
	log.Error(err)
}
