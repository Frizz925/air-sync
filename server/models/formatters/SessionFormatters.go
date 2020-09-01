package formatters

import (
	"air-sync/models"
	"air-sync/repositories/entities"
)

func SessionFromEntity(session entities.Session) models.Session {
	messages := make([]models.Message, len(session.Messages))
	for index, message := range session.Messages {
		messages[index] = MessageFromEntity(message)
	}
	return models.Session{
		ID:        session.ID,
		Messages:  messages,
		CreatedAt: session.CreatedAt,
	}
}
