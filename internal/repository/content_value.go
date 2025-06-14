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

func (r *ContentValueRepository) Create(cv *model.ContentValue) error {
	return r.db.Create(&cv).Error
}

func (r *ContentValueRepository) Save(cv *model.ContentValue) error {
	return r.db.Save(&cv).Error
}

func (r *ContentValueRepository) Delete(cv *model.ContentValue) error {
	return r.db.Delete(&cv).Error
}
