package orm

import (
	"air-sync/models"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	ID           string `gorm:"primaryKey"`
	SessionID    string `gorm:"foreignKey,not null"`
	Sensitive    bool   `gorm:"not null"`
	Body         string
	AttachmentID string
	CreatedAt    int64 `gorm:"autoCreateTime"`
}

func NewMessage(sessionID string) Message {
	return Message{
		ID:        uuid.NewV4().String(),
		SessionID: sessionID,
		Sensitive: false,
		CreatedAt: models.Timestamp(),
	}
}

func FromInsertMessageModel(sessionID string, insert models.InsertMessage) Message {
	message := NewMessage(sessionID)
	message.Sensitive = insert.Sensitive
	message.Body = insert.Body
	message.AttachmentID = insert.AttachmentID
	return message
}

func ToMessageModel(message Message) models.Message {
	return models.Message{
		BaseMessage: models.BaseMessage{
			Sensitive:    message.Sensitive,
			Body:         message.Body,
			AttachmentID: message.AttachmentID,
		},
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
	}
}
