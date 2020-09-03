package pubsub

import "sync"

type Topic struct {
	publisher   *Publisher
	name        string
	subscribers map[int]*Subscriber
	mu          sync.RWMutex
}

func NewTopic(p *Publisher, name string) *Topic {
	return &Topic{
		publisher:   p,
		name:        name,
		subscribers: make(map[int]*Subscriber),
	}
}

func (t *Topic) Subscribe() *Subscriber {
	t.mu.Lock()
	defer t.mu.Unlock()
	sub := NewSubscriber(t, t.publisher.nextID())
	t.subscribers[sub.id] = sub
	go sub.start()
	return sub
}

func (t *Topic) Publish(v interface{}) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, sub := range t.subscribers {
		sub.fire(v)
	}
}

func (t *Topic) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, sub := range t.subscribers {
		sub.cleanup()
	}
	t.subscribers = make(map[int]*Subscriber)
	t.publisher.removeTopic(t.name)
}

func (t *Topic) unsubscribe(id int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.subscribers, id)
}
