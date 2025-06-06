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
	User         *UserRepository
	Apikey       *ApikeyRepository
	Webhook      *WebhookRepository
}

func NewSet(db *gorm.DB) *Set {
	return &Set{
		Content:      NewContentRepository(db),
		Field:        NewFieldRepository(db),
		Collection:   NewCollectionRepository(db),
		ContentValue: NewContentValueRepository(db),
		Asset:        NewAssetRepository(db),
		User:         NewUserRepository(db),
		Apikey:       NewApikeyRepository(db),
		Webhook:      NewWebhookRepository(db),
	}
}
