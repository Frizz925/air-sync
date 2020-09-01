package mongo

import (
	"air-sync/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Attachment struct {
	Id        string `bson:"_id"`
	Filename  string `bson:"filename"`
	Type      string `bson:"type"`
	Mime      string `bson:"mime"`
	CreatedAt int64  `bson:"created_at"`
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
