package services

import (
	"air-sync/models/events"
	"air-sync/util/pubsub"
	"context"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack/v5"

	log "github.com/sirupsen/logrus"

	. "cloud.google.com/go/pubsub"
)

type GooglePubSubBrokerOptions struct {
	Publisher      *pubsub.Publisher
	ProjectID      string
	TopicID        string
	SubscriptionID string
}

type GooglePubSubBrokerService struct {
	context context.Context
	mu      sync.Mutex

	pub *pubsub.Publisher

	clientID       string
	projectID      string
	topicID        string
	subscriptionID string

	client *Client
	topic  *Topic
	sub    *Subscription

	lastEventID string

	initialized bool
}

var _ Initializer = (*GooglePubSubBrokerService)(nil)

func NewGooglePubSubBrokerService(ctx context.Context, opts GooglePubSubBrokerOptions) *GooglePubSubBrokerService {
	return &GooglePubSubBrokerService{
		context:        ctx,
		pub:            opts.Publisher,
		clientID:       uuid.NewV1().String(),
		projectID:      opts.ProjectID,
		topicID:        opts.TopicID,
		subscriptionID: opts.SubscriptionID,
	}
}

func (s *GooglePubSubBrokerService) Initialize() error {
	if s.initialized {
		return ErrAlreadyInitialized
	}

	log.Infof("Connecting to Google Cloud Pub/Sub: %s", s.projectID)
	client, err := NewClient(s.context, s.projectID)
	if err != nil {
		return err
	}
	s.client = client
	log.Infof("Publishing to Google Cloud Pub/Sub: %s", s.topicID)
	s.topic = client.Topic(s.topicID)
	log.Infof("Subscribing to Google Cloud Pub/Sub: %s", s.subscriptionID)
	s.sub = client.Subscription(s.subscriptionID)

	s.handlePublishingAsync()
	s.handleSubscriptionAsync()

	s.initialized = true
	return nil
}

func (s *GooglePubSubBrokerService) Deinitialize() {
	if !s.initialized {
		log.Error(ErrNotInitialized)
		return
	}
	s.initialized = false
}

func (s *GooglePubSubBrokerService) handlePublishingAsync() {
	s.pub.Topic(events.EventSession).Subscribe().
		ForEachAsync(s.context, s.handlePublishing, s.handleError)
}

func (s *GooglePubSubBrokerService) handlePublishing(v interface{}) error {
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
	b, err := msgpack.Marshal(events.PubSubSessionEvent{
		SessionEvent: event,
		ClientID:     s.clientID,
	})
	if err != nil {
		return err
	}
	res := s.topic.Publish(s.context, &Message{Data: b})
	id, err := res.Get(s.context)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"id":        id,
		"event_id":  event.ID,
		"event":     event.Event,
		"timestamp": event.Timestamp,
	}).Info("Google Cloud Pub/Sub published event")
	return nil
}

func (s *GooglePubSubBrokerService) handleSubscriptionAsync() {
	go func() {
		err := s.handleSubscription()
		if err != nil {
			s.handleError(err)
		}
	}()
}

func (s *GooglePubSubBrokerService) handleSubscription() error {
	return s.sub.Receive(s.context, func(_ context.Context, msg *Message) {
		msg.Ack()
		err := s.handleSubscriptionMessage(msg)
		if err != nil {
			s.handleError(err)
		}
	})
}

func (s *GooglePubSubBrokerService) handleSubscriptionMessage(msg *Message) error {
	event := events.PubSubSessionEvent{}
	if err := msgpack.Unmarshal(msg.Data, &event); err != nil {
		return err
	}
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
		"id":        msg.ID,
		"event_id":  event.ID,
		"event":     event.Event,
		"timestamp": event.Timestamp,
	}).Info("Google Cloud Pub/Sub received event")
	s.pub.Topic(events.EventSession).Publish(event.SessionEvent)
	return nil
}

func (s *GooglePubSubBrokerService) handleError(err error) {
	log.Error(err)
}
