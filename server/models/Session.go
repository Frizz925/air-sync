package models

import (
	"air-sync/util"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	*util.Stream
	Id string
}

func NewSession() *Session {
	return &Session{
		Stream: util.NewStream(),
		Id:     uuid.NewV4().String(),
	}
}
