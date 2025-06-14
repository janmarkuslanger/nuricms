package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type CollectionRepo interface {
	base.CRUDRepository[model.Collection]
	FindByAlias(alias string) (*model.Collection, error)
}

type collectionRepository struct {
	*base.BaseRepository[model.Collection]
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) CollectionRepo {
	return &collectionRepository{
		BaseRepository: base.NewBaseRepository[model.Collection](db),
		db:             db,
	}
}

func (r *collectionRepository) FindByAlias(alias string) (*model.Collection, error) {
	var c model.Collection
	err := r.db.Preload("Fields").Where("alias = ?", alias).First(&c).Error
	return &c, err
}

func (r *collectionRepository) FindByID(id uint, opts ...base.QueryOption) (*model.Collection, error) {
	opts = append([]base.QueryOption{base.Preload("Fields")}, opts...)
	return r.BaseRepository.FindByID(id, opts...)
}
