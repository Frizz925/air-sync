package pubsub

type Publisher interface {
	Subscribe() *Subscriber
	Unsubscribe(id int)
	ForEach(handler SubscribeFunc) error
	Fire(v interface{})
	FireError(e error)
	FireItem(item Item)
	Shutdown()
}
