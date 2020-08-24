package util

import (
	"errors"

	"github.com/reactivex/rxgo/v2"
)

var ErrStreamClosed = errors.New("Stream closed")

type Stream struct {
	rxgo.Observable
	channel chan<- rxgo.Item
}

func NewStream() *Stream {
	ch := make(chan rxgo.Item, 1)
	return &Stream{
		Observable: rxgo.FromChannel(ch),
		channel:    ch,
	}
}

func (s *Stream) Fire(event interface{}) {
	s.channel <- rxgo.Of(event)
}

func (s *Stream) FireError(err error) {
	s.channel <- rxgo.Error(err)
}

func (s *Stream) Shutdown() {
	s.FireError(ErrStreamClosed)
}
