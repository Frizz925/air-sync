package formatters

import (
	"air-sync/models"
	"air-sync/repositories/entities"
)

func MessageFromEntity(message entities.Message) models.Message {
	return models.Message{
		ID:           message.ID,
		Sensitive:    message.Sensitive,
		Body:         message.Body,
		AttachmentID: message.AttachmentID,
		CreatedAt:    message.CreatedAt,
	}
}
