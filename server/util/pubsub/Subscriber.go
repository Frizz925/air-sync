package pubsub

import (
	"errors"
	"sync"
)

var ErrSubscriberClosed = errors.New("Subscriber closed")

type SubscriberFunc func(item interface{}) error

type ErrorHandlerFunc func(err error)

type Subscriber struct {
	topic    *Topic
	id       int
	channels []chan interface{}
	mu       sync.Mutex
}

func NewSubscriber(t *Topic, id int) *Subscriber {
	return &Subscriber{
		id:       id,
		topic:    t,
		channels: make([]chan interface{}, 0),
	}
}

func (s *Subscriber) ForEach(handler SubscriberFunc) error {
	for v := range s.Channel() {
		if err := handler(v); err != nil {
			return err
		}
	}
	return nil
}

func (s *Subscriber) ForEachAsync(handler SubscriberFunc, errorHandler ErrorHandlerFunc) {
	go func() {
		if err := s.ForEach(handler); err != nil {
			errorHandler(err)
		}
	}()
}

func (s *Subscriber) Channel() <-chan interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan interface{}, 1)
	s.channels = append(s.channels, ch)
	return ch
}

func (s *Subscriber) Unsubscribe() {
	s.topic.unsubscribe(s.id)
	s.cleanup()
}

func (s *Subscriber) fire(v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.channels {
		ch <- v
	}
}

func (s *Subscriber) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.channels {
		close(ch)
	}
	s.channels = make([]chan interface{}, 0)
}
