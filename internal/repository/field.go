package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"gorm.io/gorm"
)

type FieldRepository struct {
	db *gorm.DB
}

func NewFieldRepository(db *gorm.DB) *FieldRepository {
	return &FieldRepository{db: db}
}

func (r *FieldRepository) FindByCollectionID(collectionID uint64) ([]model.Field, error) {
	var fields []model.Field
	if err := r.db.Where("collection_id = ?", collectionID).Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}
