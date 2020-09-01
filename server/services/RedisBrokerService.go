package services

import (
	"air-sync/models/events"
	"air-sync/util/pubsub"
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

type RedisBrokerOptions struct {
	Publisher *pubsub.Publisher
	Addr      string
	Password  string
}

type RedisBrokerService struct {
	context     context.Context
	pub         *pubsub.Publisher
	client      *redis.Client
	addr        string
	password    string
	initialized bool
}

func NewRedisBrokerService(opts RedisBrokerOptions) *RedisBrokerService {
	return &RedisBrokerService{
		context:     context.Background(),
		pub:         opts.Publisher,
		addr:        opts.Addr,
		password:    opts.Password,
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

	ps := client.Subscribe(s.context, events.EventSession)
	if err := ps.Ping(s.context); err != nil {
		return err
	}
	if _, err := ps.Receive(s.context); err != nil {
		return err
	}

	s.handlePublishingAsync()
	s.handleSubscriptionAsync(ps)
	log.Infof("Publishing and subscribing to Redis PubSub: %s", events.EventSession)

	s.initialized = true
	return nil
}

func (s *RedisBrokerService) Deinitialize() {
	if !s.initialized {
		log.Error(ErrNotInitialized)
		return
	}
}

func (s *RedisBrokerService) handlePublishingAsync() {
	s.pub.Topic(events.EventSession).Subscribe().
		ForEachAsync(s.handlePublishing, s.handleError)
}

func (s *RedisBrokerService) handlePublishing(v interface{}) error {
	event, ok := v.(events.SessionEvent)
	if !ok {
		return nil
	}
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return s.client.Publish(s.context, events.EventSession, string(b)).Err()
}

func (s *RedisBrokerService) handleSubscriptionAsync(ps *redis.PubSub) {
	go func() {
		err := s.handleSubscription(ps)
		if err != nil {
			s.handleError(err)
		}
	}()
}

func (s *RedisBrokerService) handleSubscription(ps *redis.PubSub) error {
	for {
		msg, err := ps.ReceiveMessage(s.context)
		if err != nil {
			return err
		}
		if msg.Channel != events.EventSession {
			continue
		}
		payload := []byte(msg.Payload)
		event := events.SessionEvent{}
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		s.pub.Topic(events.EventSession).Publish(event)
	}
}

func (s *RedisBrokerService) handleError(err error) {
	log.Error(err)
}
