package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type FieldRepo interface {
	base.CRUDRepository[model.Field]
	FindByCollectionID(collectionID uint) ([]model.Field, error)
	FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error)
}

type fieldRepository struct {
	*base.BaseRepository[model.Field]
	db *gorm.DB
}

func NewFieldRepository(db *gorm.DB) FieldRepo {
	return &fieldRepository{
		BaseRepository: base.NewBaseRepository[model.Field](db),
		db:             db,
	}
}

func (r *fieldRepository) FindByCollectionID(collectionID uint) ([]model.Field, error) {
	var fields []model.Field
	err := r.db.
		Where("collection_id = ?", collectionID).
		Find(&fields).
		Error
	return fields, err
}

func (r *fieldRepository) FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	var fields []model.Field
	err := r.db.
		Where("collection_id = ? AND display_field = ?", collectionID, true).
		Find(&fields).
		Error
	return fields, err
}
