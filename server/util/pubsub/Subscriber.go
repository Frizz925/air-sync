package pubsub

type SubscribeFunc func(v interface{}) error

type Subscriber struct {
	id        int
	publisher Publisher
	channel   chan Item
}

func (s *Subscriber) Observe() <-chan Item {
	return s.channel
}

func (s *Subscriber) Unsubscribe() {
	s.publisher.Unsubscribe(s.id)
}
