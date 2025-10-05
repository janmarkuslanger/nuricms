package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type FieldOptionRepo interface {
	base.CRUDRepository[model.FieldOption]
}

type fieldOptionRepository struct {
	*base.BaseRepository[model.FieldOption]
	db *gorm.DB
}

func NewFieldOptionRepository(db *gorm.DB) FieldOptionRepo {
	return &fieldOptionRepository{
		BaseRepository: base.NewBaseRepository[model.FieldOption](db),
		db:             db,
	}
}
