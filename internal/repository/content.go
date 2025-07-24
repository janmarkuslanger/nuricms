package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type ContentRepo interface {
	base.CRUDRepository[model.Content]
	DeleteByID(id uint) error
	FindByCollectionID(collectionID uint, offset, limit int) ([]model.Content, error)
	FindDisplayValueByCollectionID(collectionID uint, page, pageSize int) ([]model.Content, int64, error)
	ListWithDisplayContentValue() ([]model.Content, error)
	FindByCollectionAndFieldValue(collectionID uint, fieldAlias, value string, offset, limit int) ([]model.Content, int, error)
	WithTx(tx *gorm.DB) ContentRepo
}

type contentRepository struct {
	*base.BaseRepository[model.Content]
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) ContentRepo {
	return &contentRepository{
		BaseRepository: base.NewBaseRepository[model.Content](db),
		db:             db,
	}
}

func (r *contentRepository) WithTx(tx *gorm.DB) ContentRepo {
	return NewContentRepository(tx)
}

func (r *contentRepository) Create(content *model.Content) error {
	return r.BaseRepository.Create(content)
}

func (r *contentRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.Content{}, id).Error
}

func (r *contentRepository) FindByID(id uint, opts ...base.QueryOption) (*model.Content, error) {
	opts = append([]base.QueryOption{
		base.Preload("ContentValues", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Field")
		}),
		base.Preload("Collection"),
	}, opts...)
	return r.BaseRepository.FindByID(id, opts...)
}

func (r *contentRepository) FindByCollectionID(collectionID uint, offset, limit int) ([]model.Content, error) {
	db := r.db.Where("collection_id = ?", collectionID).
		Preload("ContentValues").
		Preload("ContentValues.Field")
	if offset > 0 {
		db = db.Offset(offset)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	var contents []model.Content
	return contents, db.Find(&contents).Error
}

func (r *contentRepository) FindDisplayValueByCollectionID(
	collectionID uint,
	page, pageSize int,
) ([]model.Content, int64, error) {
	var totalCount int64

	err := r.db.
		Model(&model.Content{}).
		Where("collection_id = ?", collectionID).
		Count(&totalCount).
		Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var contents []model.Content

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
		Find(&contents).
		Error

	return contents, totalCount, err
}

func (r *contentRepository) ListWithDisplayContentValue() ([]model.Content, error) {
	var contents []model.Content
	err := r.db.
		Preload("ContentValues", func(db *gorm.DB) *gorm.DB {
			return db.Joins("Field").Where("field.display_field = ?", true)
		}).
		Preload("ContentValues.Field").
		Preload("Collection").
		Find(&contents).Error
	return contents, err
}

func (r *contentRepository) FindByCollectionAndFieldValue(collectionID uint, fieldAlias, value string, offset, limit int) ([]model.Content, int, error) {
	var totalCount int64
	countDB := r.db.Model(&model.Content{}).
		Joins("JOIN content_values cv ON cv.content_id = contents.id").
		Joins("JOIN fields f ON f.id = cv.field_id").
		Where("contents.collection_id = ?", collectionID).
		Where("f.alias = ? AND cv.value = ?", fieldAlias, value)
	if err := countDB.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	var contents []model.Content
	queryDB := r.db.Model(&model.Content{}).
		Distinct("contents.*").
		Joins("JOIN content_values cv ON cv.content_id = contents.id").
		Joins("JOIN fields f ON f.id = cv.field_id").
		Where("contents.collection_id = ?", collectionID).
		Where("f.alias = ? AND cv.value = ?", fieldAlias, value)
	if offset > 0 {
		queryDB = queryDB.Offset(offset)
	}
	if limit > 0 {
		queryDB = queryDB.Limit(limit)
	}
	queryDB = queryDB.Preload("ContentValues.Field").Preload("ContentValues")
	if err := queryDB.Find(&contents).Error; err != nil {
		return nil, 0, err
	}
	return contents, int(totalCount), nil
}
