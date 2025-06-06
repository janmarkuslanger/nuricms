package repository

import (
	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

type FieldRepository struct {
	db *gorm.DB
}

func NewFieldRepository(db *gorm.DB) *FieldRepository {
	return &FieldRepository{db: db}
}

func (r *FieldRepository) FindByCollectionID(collectionID uint) ([]model.Field, error) {
	var fields []model.Field
	if err := r.db.Where("collection_id = ?", collectionID).Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

func (r *FieldRepository) FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	var fields []model.Field

	if err := r.db.Where("collection_id = ? AND display_field = 1", collectionID).Find(&fields).Error; err != nil {
		return nil, err
	}

	return fields, nil
}

func (r *FieldRepository) List(page, pageSize int) ([]model.Field, int64, error) {
	var fields []model.Field
	var totalCount int64

	err := r.db.Model(&model.Field{}).
		Preload("Collection").
		Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err = r.db.
		Preload("Collection").
		Offset(offset).
		Limit(pageSize).
		Find(&fields).Error

	return fields, totalCount, err
}
