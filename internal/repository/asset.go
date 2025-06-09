package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
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

func (r *AssetRepository) List(page, pageSize int) ([]model.Asset, int64, error) {
	var assets []model.Asset
	var totalCount int64

	if err := r.db.Model(&model.Asset{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := r.db.Offset(offset).Limit(pageSize).Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, totalCount, nil
}

func (r *AssetRepository) Save(asset *model.Asset) error {
	if err := r.db.Save(asset).Error; err != nil {
		return err
	}

	return nil
}
