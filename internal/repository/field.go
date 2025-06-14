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

func (r *FieldRepository) FindByID(id uint) (*model.Field, error) {
	var field model.Field
	if err := r.db.Preload("Collection").First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

func (r *FieldRepository) Save(field *model.Field) error {
	return r.db.Save(field).Error
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

func (r *FieldRepository) Create(field *model.Field) error {
	return r.db.Create(field).Error
}

func (r *FieldRepository) Delete(field *model.Field) error {
	return r.db.Delete(field).Error
}
