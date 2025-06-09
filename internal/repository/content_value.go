package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"gorm.io/gorm"
)

type ContentValueRepository struct {
	db *gorm.DB
}

func NewContentValueRepository(db *gorm.DB) *ContentValueRepository {
	return &ContentValueRepository{db: db}
}

func (r *ContentValueRepository) Create(contentValue *model.ContentValue) error {
	result := r.db.Create(&contentValue)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
