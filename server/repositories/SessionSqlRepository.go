package repositories

import (
	"air-sync/models"
	"air-sync/repositories/entities"
	"errors"

	"gorm.io/gorm"
)

type SessionSqlRepository struct {
	*SqlRepository
}

var _ SessionRepository = (*SessionSqlRepository)(nil)
var _ RepositoryMigration = (*SessionSqlRepository)(nil)

func NewSessionSqlRepository(db *gorm.DB) *SessionSqlRepository {
	return &SessionSqlRepository{NewSqlRepository(db)}
}

func (r *SessionSqlRepository) Migrate() error {
	if err := r.db.AutoMigrate(entities.Session{}); err != nil {
		return err
	}
	if err := r.db.AutoMigrate(entities.Message{}); err != nil {
		return err
	}
	return nil
}

func (r *SessionSqlRepository) Create() (entities.Session, error) {
	session := entities.NewSession()
	err := r.db.Create(&session).Error
	return session, r.sessionCrudError(err)
}

func (r *SessionSqlRepository) Find(id string) (entities.Session, error) {
	session := entities.Session{}
	err := r.db.First(&session, id).Error
	return session, r.sessionCrudError(err)
}

func (r *SessionSqlRepository) InsertMessage(id string, model models.Message) (entities.Message, error) {
	message := entities.FromMessageModel(id, model)
	err := r.db.Create(&message).Error
	return message, r.messageCrudError(err)
}

func (r *SessionSqlRepository) DeleteMessage(id string, messageId string) error {
	err := r.db.Delete(entities.Message{
		ID:        messageId,
		SessionID: id,
	}).Error
	return r.messageCrudError(err)
}

func (r *SessionSqlRepository) Delete(id string) error {
	err := r.db.Delete(entities.Session{}, id).Error
	return r.sessionCrudError(err)
}

func (r *SessionSqlRepository) sessionCrudError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrSessionNotFound
	}
	return err
}

func (r *SessionSqlRepository) messageCrudError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrMessageNotFound
	}
	return err
}
