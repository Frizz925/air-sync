package mongo

import (
	"air-sync/models"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	ID           string `bson:"id"`
	SessionID    string `bson:"session_id"`
	Sensitive    bool   `bson:"sensitive"`
	Body         string `bson:"body"`
	AttachmentID string `bson:"attachment_id"`
	CreatedAt    int64  `bson:"created_at"`
}

func NewMessage(sessionId string) Message {
	return Message{
		ID:        uuid.NewV4().String(),
		SessionID: sessionId,
		Sensitive: false,
		CreatedAt: models.Timestamp(),
	}
}

func FromInsertMessageModel(sessionId string, insert models.InsertMessage) Message {
	message := NewMessage(sessionId)
	message.Sensitive = insert.Sensitive
	message.Body = insert.Body
	message.AttachmentID = insert.AttachmentID
	return message
}

func ToMessageModel(message Message, attachment Attachment) models.Message {
	return models.Message{
		BaseMessage: models.BaseMessage{
			Sensitive:    message.Sensitive,
			Body:         message.Body,
			AttachmentID: message.AttachmentID,
		},
		ID:             message.ID,
		AttachmentType: attachment.Type,
		AttachmentName: attachment.Name,
		CreatedAt:      message.CreatedAt,
	}
}
