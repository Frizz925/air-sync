package pubsub

import "sync"

type Message struct {
	Topic string
	Value interface{}
}

type Publisher struct {
	topics  map[string]*Topic
	_nextID int
	mu      sync.RWMutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		_nextID: 0,
		topics:  make(map[string]*Topic),
	}
}

func (p *Publisher) Topic(name string) *Topic {
	p.mu.Lock()
	defer p.mu.Unlock()
	if topic, ok := p.topics[name]; ok {
		return topic
	}
	topic := NewTopic(p, name)
	p.topics[name] = topic
	return topic
}

func (p *Publisher) nextID() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p._nextID++
	return p._nextID
}

func (p *Publisher) removeTopic(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.topics, name)
}
