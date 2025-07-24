package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type ContentValueRepo interface {
	base.CRUDRepository[model.ContentValue]
	FindByContentID(cID uint) ([]model.ContentValue, error)
	WithTx(tx *gorm.DB) ContentValueRepo
}

type contentValueRepository struct {
	*base.BaseRepository[model.ContentValue]
	db *gorm.DB
}

func NewContentValueRepository(db *gorm.DB) ContentValueRepo {
	return &contentValueRepository{
		BaseRepository: base.NewBaseRepository[model.ContentValue](db),
		db:             db,
	}
}

func (r *contentValueRepository) WithTx(tx *gorm.DB) ContentValueRepo {
	return NewContentValueRepository(tx)
}

func (r *contentValueRepository) FindByContentID(cID uint) ([]model.ContentValue, error) {
	var cvs []model.ContentValue
	err := r.db.
		Where("content_id = ?", cID).
		Find(&cvs).
		Error
	return cvs, err
}
