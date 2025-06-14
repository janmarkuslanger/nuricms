package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"gorm.io/gorm"
)

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) Create(content *model.Content) error {
	return r.db.Create(&content).Error
}

func (r *ContentRepository) FindByID(id uint) (model.Content, error) {
	var content model.Content
	err := r.db.
		Where("id = ?", id).
		Preload("ContentValues").
		Preload("ContentValues.Field").
		Preload("Collection").
		Find(&content).Error
	return content, err
}

func (r *ContentRepository) FindByCollectionID(collectionID uint, offset int, limit int) ([]model.Content, error) {
	q := r.db.
		Where("collection_id = ?", collectionID).
		Preload("ContentValues").
		Preload("ContentValues.Field")

	if offset > 0 {
		q = q.Offset(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	var contents []model.Content
	if err := q.Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}

func (r *ContentRepository) FindDisplayValueByCollectionID(collectionID uint, page, pageSize int) ([]model.Content, int64, error) {
	var contents []model.Content
	var totalCount int64

	err := r.db.Model(&model.Content{}).
		Where("collection_id = ?", collectionID).
		Preload("ContentValues", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("Field").
				Where("field.display_field = ?", true)
		}).
		Preload("ContentValues.Field").
		Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err = r.db.
		Where("collection_id = ?", collectionID).
		Preload("ContentValues", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("Field").
				Where("field.display_field = ?", true)
		}).
		Preload("ContentValues.Field").
		Offset(offset).
		Limit(pageSize).
		Find(&contents).Error

	return contents, totalCount, err
}

func (r *ContentRepository) ListWithDisplayContentValue() ([]model.Content, error) {
	var contents []model.Content

	err := r.db.
		Preload("ContentValues", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("Field").
				Where("field.display_field = ?", true)
		}).
		Preload("ContentValues.Field").
		Preload("Collection").
		Find(&contents).Error

	return contents, err
}

func (r *ContentRepository) FindByCollectionAndFieldValue(collectionID uint, fieldAlias, value string, offset, limit int) ([]model.Content, int, error) {
	var contents []model.Content
	var totalCount int64

	countErr := r.db.
		Model(&model.Content{}).
		Joins("JOIN content_values cv ON cv.content_id = contents.id").
		Joins("JOIN fields f ON f.id = cv.field_id").
		Where("contents.collection_id = ?", collectionID).
		Where("f.alias = ? AND cv.value = ?", fieldAlias, value).
		Count(&totalCount).Error
	if countErr != nil {
		return nil, 0, countErr
	}
	if totalCount == 0 {
		return []model.Content{}, 0, nil
	}

	queryErr := r.db.
		Model(&model.Content{}).
		Distinct("contents.*").
		Preload("ContentValues.Field").
		Preload("ContentValues").
		Joins("JOIN content_values cv ON cv.content_id = contents.id").
		Joins("JOIN fields f ON f.id = cv.field_id").
		Where("contents.collection_id = ?", collectionID).
		Where("f.alias = ? AND cv.value = ?", fieldAlias, value).
		Offset(offset).
		Limit(limit).
		Find(&contents).Error
	if queryErr != nil {
		return nil, 0, queryErr
	}
	return contents, int(totalCount), nil
}
