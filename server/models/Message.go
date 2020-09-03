package models

type BaseMessage struct {
	Sensitive      bool   `json:"sensitive"`
	Body           string `json:"body,omitempty"`
	AttachmentID   string `json:"attachment_id,omitempty"`
	AttachmentName string `json:"attachment_name,omitempty"`
}

type Message struct {
	BaseMessage
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
}

type InsertMessage struct {
	BaseMessage
}

var EmptyMessage = Message{}
