package models

type BaseMessage struct {
	Sensitive    bool   `json:"sensitive"`
	Body         string `json:"body,omitempty"`
	AttachmentId string `json:"attachment_id,omitempty"`
}

type Message struct {
	BaseMessage
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
}

type InsertMessage struct {
	BaseMessage
}

var EmptyMessage = Message{}
