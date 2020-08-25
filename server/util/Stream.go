package util

import (
	"context"
	"errors"

	"github.com/reactivex/rxgo/v2"
)

var ErrStreamClosed = errors.New("Stream closed")

type Stream struct {
	context.Context
	rxgo.Observable
	channel chan<- rxgo.Item
	dispose rxgo.Disposable
}

func NewStream() *Stream {
	ch := make(chan rxgo.Item, 1)
	ob := rxgo.FromChannel(ch, rxgo.WithPublishStrategy(), rxgo.WithBufferedChannel(1))
	ctx, dispose := ob.Connect()
	return &Stream{
		Context:    ctx,
		Observable: ob,
		channel:    ch,
		dispose:    dispose,
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
	close(s.channel)
	s.dispose()
}
