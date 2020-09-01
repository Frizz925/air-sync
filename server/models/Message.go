package models

type Message struct {
	ID           string `json:"id"`
	Sensitive    bool   `json:"sensitive"`
	Body         string `json:"body"`
	AttachmentID string `json:"attachment_id"`
	CreatedAt    int64  `json:"created_at"`
}

var EmptyMessage = Message{}
