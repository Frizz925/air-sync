package repositories

import (
	"air-sync/models"
	"air-sync/models/orm"
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
	if err := r.db.AutoMigrate(orm.Session{}); err != nil {
		return err
	}
	if err := r.db.AutoMigrate(orm.Message{}); err != nil {
		return err
	}
	return nil
}

func (r *SessionSqlRepository) Create() (models.Session, error) {
	session := orm.NewSession()
	err := r.db.Create(&session).Error
	return orm.ToSessionModel(session), r.sessionCrudError(err)
}

func (r *SessionSqlRepository) Find(id string) (models.Session, error) {
	session := orm.Session{}
	err := r.db.First(&session, id).Error
	return orm.ToSessionModel(session), r.sessionCrudError(err)
}

func (r *SessionSqlRepository) InsertMessage(id string, arg models.InsertMessage) (models.Message, error) {
	message := orm.FromInsertMessageModel(id, arg)
	err := r.db.Create(&message).Error
	return orm.ToMessageModel(message), r.messageCrudError(err)
}

func (r *SessionSqlRepository) DeleteMessage(id string, messageID string) error {
	err := r.db.Delete(orm.Message{
		ID:        messageID,
		SessionID: id,
	}).Error
	return r.messageCrudError(err)
}

func (r *SessionSqlRepository) Delete(id string) error {
	err := r.db.Delete(orm.Session{}, id).Error
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
