package pubsub

import (
	"context"
	"errors"
	"sync"
)

var ErrSubscriberClosed = errors.New("Subscriber closed")

type SubscriberFunc func(item interface{}) error

type ErrorHandlerFunc func(err error)

type Subscriber struct {
	context   context.Context
	topic     *Topic
	id        int
	sender    chan interface{}
	receivers []chan interface{}
	cancel    context.CancelFunc
	mu        sync.Mutex
}

func NewSubscriber(t *Topic, id int) *Subscriber {
	ctx, cancel := context.WithCancel(context.Background())
	return &Subscriber{
		context:   ctx,
		topic:     t,
		id:        id,
		sender:    make(chan interface{}, 1),
		receivers: make([]chan interface{}, 0),
		cancel:    cancel,
	}
}

func (s *Subscriber) ForEach(ctx context.Context, handler SubscriberFunc) error {
	ch := s.Channel()
	for {
		select {
		case v := <-ch:
			if err := handler(v); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Subscriber) ForEachAsync(ctx context.Context, handler SubscriberFunc, errorHandler ErrorHandlerFunc) {
	go func() {
		if err := s.ForEach(ctx, handler); err != nil {
			errorHandler(err)
		}
	}()
}

func (s *Subscriber) Channel() <-chan interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan interface{}, 1)
	s.receivers = append(s.receivers, ch)
	return ch
}

func (s *Subscriber) Unsubscribe() {
	s.topic.unsubscribe(s.id)
	s.cleanup()
}

func (s *Subscriber) start() {
	for {
		select {
		case v := <-s.sender:
			s.mu.Lock()
			for _, ch := range s.receivers {
				ch <- v
			}
			s.mu.Unlock()
		case <-s.context.Done():
			return
		}
	}
}

func (s *Subscriber) fire(v interface{}) {
	s.sender <- v
}

func (s *Subscriber) cleanup() {
	s.cancel()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.receivers {
		close(ch)
	}
	s.receivers = make([]chan interface{}, 0)
}
