package repository

import (
	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) Create(content *model.Content) (model.Content, error) {
	if err := r.db.Create(&content).Error; err != nil {
		return *content, err
	}
	return *content, nil
}

func (r *ContentRepository) FindByCollectionID(collectionID uint) ([]model.Content, error) {
	var contents []model.Content
	err := r.db.
		Where("collection_id = ?", collectionID).
		Preload("ContentValues").
		Preload("ContentValues.Field").
		Find(&contents).Error
	return contents, err
}

func (r *ContentRepository) FindDisplayValueByCollectionID(collectionID uint) ([]model.Content, error) {
	var contents []model.Content

	err := r.db.
		Where("collection_id = ?", collectionID).
		Preload("ContentValues", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("Field").
				Where("field.display_field = ?", true)
		}).
		Preload("ContentValues.Field").
		Find(&contents).Error

	return contents, err
}
