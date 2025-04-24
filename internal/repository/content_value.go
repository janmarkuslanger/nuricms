package repository

import "gorm.io/gorm"

type ContentValueRepository struct {
	db *gorm.DB
}

func NewContentValueRepository(db *gorm.DB) *ContentValueRepository {
	return &ContentValueRepository{db: db}
}
