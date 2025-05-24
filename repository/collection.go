package repository

import (
	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) FindByID(id uint) (*model.Collection, error) {
	var collection model.Collection
	if err := r.db.Preload("Fields").First(&collection, id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *CollectionRepository) FindByAlias(alias string) (*model.Collection, error) {
	var collection model.Collection

	err := r.db.
		Preload("Fields").
		Where("alias = ?", alias).
		First(&collection).Error
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (r *CollectionRepository) List(page, pageSize int) ([]model.Collection, int64, error) {
	var collections []model.Collection
	var totalCount int64

	if err := r.db.Model(&model.Collection{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := r.db.Offset(offset).Limit(pageSize).Find(&collections).Error; err != nil {
		return nil, 0, err
	}

	return collections, totalCount, nil
}
