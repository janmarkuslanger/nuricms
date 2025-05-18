// repository/api_key_repository.go
package repository

import (
	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

type ApikeyRepository struct {
	db *gorm.DB
}

func NewApikeyRepository(db *gorm.DB) *ApikeyRepository {
	return &ApikeyRepository{db: db}
}

func (r *ApikeyRepository) Create(key *model.Apikey) error {
	return r.db.Create(key).Error
}

func (r *ApikeyRepository) List() ([]model.Apikey, error) {
	var keys []model.Apikey
	if err := r.db.Find(&keys).Error; err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *ApikeyRepository) Delete(apikey *model.Apikey) error {
	err := r.db.Delete(apikey).Error
	return err
}

func (r *ApikeyRepository) FindByToken(token string) (*model.Apikey, error) {
	var apikey model.Apikey
	if err := r.db.Where("token = ?", token).First(&apikey).Error; err != nil {
		return nil, err
	}
	return &apikey, nil
}
