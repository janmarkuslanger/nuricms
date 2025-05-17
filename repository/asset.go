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

func (r *AssetRepository) Create(asset *model.Asset) error {
	if err := r.db.Create(&asset).Error; err != nil {
		return err
	}
	return nil
}

func (r *AssetRepository) Delete(asset *model.Asset) error {
	err := r.db.Delete(asset).Error
	return err
}

func (r *AssetRepository) FindByID(id uint) (*model.Asset, error) {
	var asset model.Asset

	if err := r.db.First(&asset, id).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *AssetRepository) List() ([]model.Asset, error) {
	var assets []model.Asset
	if err := r.db.Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *AssetRepository) Save(asset *model.Asset) error {
	if err := r.db.Save(asset).Error; err != nil {
		return err
	}

	return nil
}
