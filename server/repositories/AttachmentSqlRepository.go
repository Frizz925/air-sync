package repositories

import (
	"air-sync/repositories/entities"
	"errors"

	"gorm.io/gorm"
)

type AttachmentSqlRepository struct {
	*SqlRepository
}

var _ AttachmentRepository = (*AttachmentSqlRepository)(nil)
var _ RepositoryMigration = (*AttachmentSqlRepository)(nil)

func NewAttachmentSqlRepository(db *gorm.DB) *AttachmentSqlRepository {
	return &AttachmentSqlRepository{NewSqlRepository(db)}
}

func (r *AttachmentSqlRepository) Migrate() error {
	return r.db.AutoMigrate(entities.Attachment{})
}

func (r *AttachmentSqlRepository) Create(arg entities.CreateAttachment) (entities.Attachment, error) {
	attachment := entities.FromCreateAttachment(arg)
	err := r.db.Create(&attachment).Error
	return attachment, r.crudError(err)
}

func (r *AttachmentSqlRepository) Find(id string) (entities.Attachment, error) {
	attachment := entities.Attachment{}
	err := r.db.First(&attachment).Error
	return attachment, r.crudError(err)
}

func (r *AttachmentSqlRepository) Delete(id string) error {
	err := r.db.Delete(entities.Attachment{}, id).Error
	return r.crudError(err)
}

func (r *AttachmentSqlRepository) crudError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAttachmentNotFound
	}
	return err
}
