package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type AssetRepo interface {
	base.CRUDRepository[model.Asset]
}

type AssetRepository struct {
	*base.BaseRepository[model.Asset]
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{
		BaseRepository: base.NewBaseRepository[model.Asset](db),
		db:             db,
	}
}
