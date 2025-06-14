package repository

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type UserRepo interface {
	base.CRUDRepository[model.User]
	FindByEmail(email string) (*model.User, error)
}

type userRepository struct {
	*base.BaseRepository[model.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &userRepository{
		BaseRepository: base.NewBaseRepository[model.User](db),
		db:             db,
	}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &u, err
}
