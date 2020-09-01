package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Attachment struct {
	ID        string `gorm:"primaryKey"`
	Filename  string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Mime      string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}

type CreateAttachment struct {
	Filename string
	Type     string
	Mime     string
}

func FromCreateAttachment(attachment CreateAttachment) Attachment {
	return Attachment{
		ID:        uuid.NewV4().String(),
		Filename:  attachment.Filename,
		Type:      attachment.Type,
		Mime:      attachment.Mime,
		CreatedAt: time.Now().Unix(),
	}
}
