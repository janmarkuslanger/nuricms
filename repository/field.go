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
