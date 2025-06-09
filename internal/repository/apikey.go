package repository

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
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

func (r *ApikeyRepository) List(page, pageSize int) ([]model.Apikey, int64, error) {
	var keys []model.Apikey
	var totalCount int64

	if err := r.db.Model(&model.Apikey{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := r.db.Offset(offset).Limit(pageSize).Find(&keys).Error; err != nil {
		return nil, 0, err
	}

	return keys, totalCount, nil
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

func (r *ApikeyRepository) FindByID(id uint) (*model.Apikey, error) {
	var apikey model.Apikey
	if err := r.db.Where("id = ?", id).First(&apikey).Error; err != nil {
		return nil, err
	}
	return &apikey, nil
}
