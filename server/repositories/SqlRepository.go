package repositories

import (
	"gorm.io/gorm"
)

type SqlRepository struct {
	db *gorm.DB
}

func NewSqlRepository(db *gorm.DB) *SqlRepository {
	return &SqlRepository{db}
}
