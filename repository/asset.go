package repository

import (
	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(asset *model.Asset) (model.Asset, error) {
	if err := r.db.Create(&asset).Error; err != nil {
		return *asset, err
	}
	return *asset, nil
}
