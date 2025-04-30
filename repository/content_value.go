package repository

import (
	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

type ContentValueRepository struct {
	db *gorm.DB
}

func NewContentValueRepository(db *gorm.DB) *ContentValueRepository {
	return &ContentValueRepository{db: db}
}

func (r *ContentValueRepository) Create(contentValue *model.ContentValue) (model.ContentValue, error) {
	result := r.db.Create(&contentValue)

	if result.Error != nil {
		return *contentValue, result.Error
	}

	return *contentValue, nil
}
