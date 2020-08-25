package models

import (
	uuid "github.com/satori/go.uuid"
)

type Content struct {
	Type    string `json:"type"`
	Mime    string `json:"mime"`
	Payload string `json:"payload"`
}

type Session struct {
	Id      string   `json:"id"`
	Content *Content `json:"content"`
}

func NewSession() *Session {
	return &Session{
		Id: uuid.NewV4().String(),
		Content: &Content{
			Type:    "text",
			Mime:    "text/plain",
			Payload: "",
		},
	}
}
