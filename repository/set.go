package repository

import (
	"gorm.io/gorm"
)

type Set struct {
	Content      *ContentRepository
	Field        *FieldRepository
	Collection   *CollectionRepository
	ContentValue *ContentValueRepository
	Asset        *AssetRepository
}

func NewSet(db *gorm.DB) *Set {
	return &Set{
		Content:      NewContentRepository(db),
		Field:        NewFieldRepository(db),
		Collection:   NewCollectionRepository(db),
		ContentValue: NewContentValueRepository(db),
		Asset:        NewAssetRepository(db),
	}
}
