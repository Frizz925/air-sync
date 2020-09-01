package orm

import (
	"air-sync/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Attachment struct {
	Id        string `gorm:"primaryKey"`
	Filename  string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Mime      string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}

func NewAttachment() Attachment {
	return Attachment{
		Id:        uuid.NewV4().String(),
		CreatedAt: time.Now().Unix(),
	}
}

func FromCreateAttachmentModel(create models.CreateAttachment) Attachment {
	attachment := NewAttachment()
	attachment.Type = create.Type
	attachment.Filename = create.Filename
	attachment.Mime = create.Mime
	return attachment
}

func ToAttachmentModel(attachment Attachment) models.Attachment {
	return models.Attachment{
		BaseAttachment: models.BaseAttachment{
			Type:     attachment.Type,
			Filename: attachment.Filename,
			Mime:     attachment.Mime,
		},
		Id:        attachment.Id,
		CreatedAt: attachment.CreatedAt,
	}
}
