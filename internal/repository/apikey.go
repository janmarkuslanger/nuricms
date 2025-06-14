package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type ApikeyRepo interface {
	base.CRUDRepository[model.Apikey]
	FindByToken(token string) (*model.Apikey, error)
}

type apikeyRepository struct {
	*base.BaseRepository[model.Apikey]
	db *gorm.DB
}

func NewApikeyRepository(db *gorm.DB) *apikeyRepository {
	return &apikeyRepository{
		BaseRepository: base.NewBaseRepository[model.Apikey](db),
		db:             db,
	}
}

func (r *apikeyRepository) FindByToken(token string) (*model.Apikey, error) {
	var ak model.Apikey
	err := r.db.Where("token = ?", token).First(&ak).Error
	return &ak, err
}
