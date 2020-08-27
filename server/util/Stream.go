package util

import (
	"errors"
	"sync"
)

var ErrStreamClosed = errors.New("Stream closed")

type Item struct {
	E error
	V interface{}
}

type Subscriber struct {
	id      int
	stream  *Stream
	channel chan *Item
}

type Stream struct {
	sync.RWMutex
	nextId      int
	subscribers map[int]*Subscriber
}

func NewStream() *Stream {
	return &Stream{
		nextId:      0,
		subscribers: make(map[int]*Subscriber),
	}
}

func (s *Stream) Subscribe() *Subscriber {
	defer s.Unlock()
	s.Lock()
	s.nextId++
	sub := &Subscriber{
		id:      s.nextId,
		stream:  s,
		channel: make(chan *Item, 1),
	}
	s.subscribers[sub.id] = sub
	return sub
}

func (s *Stream) Unsubscribe(id int) {
	defer s.Unlock()
	s.Lock()
	if sub, ok := s.subscribers[id]; ok {
		delete(s.subscribers, id)
		close(sub.channel)
	}
}

func (s *Stream) Fire(event interface{}) {
	s.FireItem(&Item{
		E: nil,
		V: event,
	})
}

func (s *Stream) FireError(err error) {
	s.FireItem(&Item{
		E: err,
		V: nil,
	})
}

func (s *Stream) FireItem(item *Item) {
	defer s.RUnlock()
	s.RLock()
	for _, sub := range s.subscribers {
		select {
		case sub.channel <- item:
		default:
		}
	}
}

func (s *Stream) Shutdown() {
	s.FireError(ErrStreamClosed)
	for _, sub := range s.subscribers {
		s.Unsubscribe(sub.id)
	}
}

func (s *Subscriber) Observe() <-chan *Item {
	return s.channel
}

func (s *Subscriber) Unsubscribe() {
	s.stream.Unsubscribe(s.id)
}
