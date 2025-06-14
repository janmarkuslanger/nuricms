package repository

import (
	"gorm.io/gorm"
)

type Set struct {
	Content      *ContentRepository
	Field        *FieldRepository
	Collection   CollectionRepo
	ContentValue *ContentValueRepository
	Asset        AssetRepo
	User         *UserRepository
	Apikey       ApikeyRepo
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
