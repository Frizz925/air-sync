package mongo

import (
	"air-sync/models"

	uuid "github.com/satori/go.uuid"
)

type Attachment struct {
	ID        string `bson:"_id"`
	Type      string `bson:"type"`
	Mime      string `bson:"mime"`
	Name      string `bson:"name"`
	CreatedAt int64  `bson:"created_at"`
}

var EmptyAttachment = Attachment{}

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
