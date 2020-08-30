package pubsub

import "sync"

type Topic struct {
	sync.RWMutex
	stream      *Stream
	name        string
	nextId      int
	subscribers map[int]*Subscriber
}

var _ Publisher = (*Topic)(nil)

func (t *Topic) Subscribe() *Subscriber {
	defer t.Unlock()
	t.Lock()
	t.nextId++
	sub := &Subscriber{
		id:        t.nextId,
		publisher: t,
		channel:   make(chan Item, 1),
	}
	t.subscribers[sub.id] = sub
	return sub
}

func (t *Topic) Unsubscribe(id int) {
	defer t.Unlock()
	t.Lock()
	if sub, ok := t.subscribers[id]; ok {
		delete(t.subscribers, id)
		close(sub.channel)
	}
}

func (t *Topic) Fire(event interface{}) {
	t.FireItem(Item{
		E: nil,
		V: event,
	})
}

func (t *Topic) FireError(err error) {
	t.FireItem(Item{
		E: err,
		V: nil,
	})
}

func (t *Topic) FireItem(item Item) {
	defer t.RUnlock()
	t.RLock()
	for _, sub := range t.subscribers {
		select {
		case sub.channel <- item:
		default:
		}
	}
}

func (t *Topic) ForEach(handler SubscribeFunc) error {
	sub := t.Subscribe()
	defer sub.Unsubscribe()
	for item := range sub.Observe() {
		if item.E == ErrStreamClosed {
			break
		}
		if err := handler(item.V); err != nil {
			return err
		}
	}
	return nil
}

func (t *Topic) Shutdown() {
	t.FireError(ErrStreamClosed)
	for _, sub := range t.subscribers {
		t.Unsubscribe(sub.id)
	}
}
