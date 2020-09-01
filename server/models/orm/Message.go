package orm

import (
	"air-sync/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	Id           string `gorm:"primaryKey"`
	SessionId    string `gorm:"foreignKey,not null"`
	Sensitive    bool   `gorm:"not null"`
	Body         string
	AttachmentId string
	CreatedAt    int64 `gorm:"autoCreateTime"`
}

func NewMessage(sessionId string) Message {
	return Message{
		Id:        uuid.NewV4().String(),
		SessionId: sessionId,
		Sensitive: false,
		CreatedAt: time.Now().Unix(),
	}
}

func FromInsertMessageModel(sessionId string, insert models.InsertMessage) Message {
	message := NewMessage(sessionId)
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
