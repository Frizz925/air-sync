package orm

import (
	"air-sync/models"

	uuid "github.com/satori/go.uuid"
)

type Attachment struct {
	ID        string `gorm:"primaryKey"`
	Type      string `gorm:"not null"`
	Mime      string `gorm:"not null"`
	Name      string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}

func NewAttachment() Attachment {
	return Attachment{
		ID:        uuid.NewV4().String(),
		CreatedAt: models.Timestamp(),
	}
}

func FromCreateAttachmentModel(create models.CreateAttachment) Attachment {
	attachment := NewAttachment()
	attachment.Type = create.Type
	attachment.Mime = create.Mime
	attachment.Name = create.Name
	return attachment
}

func ToAttachmentModel(attachment Attachment) models.Attachment {
	return models.Attachment{
		BaseAttachment: models.BaseAttachment{
			Type: attachment.Type,
			Mime: attachment.Mime,
			Name: attachment.Name,
		},
		ID:        attachment.ID,
		CreatedAt: attachment.CreatedAt,
	}
}
