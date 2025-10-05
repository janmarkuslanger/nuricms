package repository

import (
	"gorm.io/gorm"
)

type Set struct {
	Content      ContentRepo
	Field        FieldRepo
	FieldOption  FieldOptionRepo
	Collection   CollectionRepo
	ContentValue ContentValueRepo
	Asset        AssetRepo
	User         UserRepo
	Apikey       ApikeyRepo
	Webhook      WebhookRepo
}

func NewSet(db *gorm.DB) *Set {
	return &Set{
		Content:      NewContentRepository(db),
		Field:        NewFieldRepository(db),
		FieldOption:  NewFieldOptionRepository(db),
		Collection:   NewCollectionRepository(db),
		ContentValue: NewContentValueRepository(db),
		Asset:        NewAssetRepository(db),
		User:         NewUserRepository(db),
		Apikey:       NewApikeyRepository(db),
		Webhook:      NewWebhookRepository(db),
	}
}
