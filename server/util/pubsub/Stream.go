package pubsub

import (
	"errors"
	"sync"
)

var ErrStreamClosed = errors.New("Stream closed")

type Stream struct {
	sync.Mutex
	topics map[string]*Topic
}

func NewStream() *Stream {
	return &Stream{
		topics: make(map[string]*Topic),
	}
}

func (s *Stream) Topic(name string) *Topic {
	defer s.Unlock()
	s.Lock()
	topic := s.topics[name]
	if topic == nil {
		topic = &Topic{
			stream:      s,
			name:        name,
			nextId:      0,
			subscribers: make(map[int]*Subscriber),
		}
		s.topics[name] = topic
	}
	return topic
}

func (s *Stream) DeleteTopic(name string) {
	defer s.Unlock()
	s.Lock()
	delete(s.topics, name)
}

func (s *Stream) Shutdown() {
	defer s.Unlock()
	s.Lock()
	for name, topic := range s.topics {
		topic.Shutdown()
		delete(s.topics, name)
	}
}
