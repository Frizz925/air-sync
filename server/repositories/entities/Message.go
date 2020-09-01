package entities

import (
	"air-sync/models"
	"time"

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

func FromMessageModel(sessionId string, message models.Message) Message {
	return Message{
		ID:           uuid.NewV4().String(),
		SessionID:    sessionId,
		Sensitive:    message.Sensitive,
		Body:         message.Body,
		AttachmentID: message.AttachmentID,
		CreatedAt:    time.Now().Unix(),
	}
}
