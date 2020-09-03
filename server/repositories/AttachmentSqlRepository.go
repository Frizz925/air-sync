package repositories

import (
	"air-sync/models"
	"air-sync/models/orm"
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
	return r.db.AutoMigrate(orm.Attachment{})
}

func (r *AttachmentSqlRepository) Create(arg models.CreateAttachment) (models.Attachment, error) {
	attachment := orm.FromCreateAttachmentModel(arg)
	err := r.db.Create(&attachment).Error
	return orm.ToAttachmentModel(attachment), r.crudError(err)
}

func (r *AttachmentSqlRepository) Find(id string) (models.Attachment, error) {
	attachment := orm.Attachment{}
	err := r.db.First(&attachment, id).Error
	return orm.ToAttachmentModel(attachment), r.crudError(err)
}

func (r *AttachmentSqlRepository) Delete(id string) error {
	err := r.db.Delete(orm.Attachment{}, id).Error
	return r.crudError(err)
}

func (r *AttachmentSqlRepository) crudError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAttachmentNotFound
	}
	return err
}
