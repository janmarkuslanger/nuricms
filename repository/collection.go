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

func (r *CollectionRepository) FindByID(id uint64) (*model.Collection, error) {
	var collection model.Collection
	if err := r.db.Preload("Fields").First(&collection, id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *CollectionRepository) GetAll() ([]model.Collection, error) {
	var collections []model.Collection
	if err := r.db.Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}
