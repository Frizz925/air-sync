package mongo

import (
	"air-sync/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	Id           string `bson:"_id"`
	Sensitive    bool   `bson:"sensitive"`
	Body         string `bson:"body"`
	AttachmentId string `bson:"attachment_id"`
	CreatedAt    int64  `bson:"created_at"`
}

func NewMessage() Message {
	return Message{
		Id:        uuid.NewV4().String(),
		Sensitive: false,
		CreatedAt: time.Now().Unix(),
	}
}

func FromInsertMessageModel(insert models.InsertMessage) Message {
	message := NewMessage()
	message.Sensitive = insert.Sensitive
	message.Body = insert.Body
	message.AttachmentId = insert.AttachmentId
	return message
}

func ToMessageModel(message Message) models.Message {
	return models.Message{
		BaseMessage: models.BaseMessage{
			Sensitive:    message.Sensitive,
			Body:         message.Body,
			AttachmentId: message.AttachmentId,
		},
		Id:        message.Id,
		CreatedAt: message.CreatedAt,
	}
}
