package service

import (
	"log"
	"os"

	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
)

type Set struct {
	Collection   *CollectionService
	Field        *FieldService
	Content      *ContentService
	ContentValue *ContentValueService
	Asset        *AssetService
	User         *UserService
	Apikey       *ApikeyService
	Webhook      *WebhookService
	Api          *ApiService
}

func NewSet(r *repository.Set, hr *plugin.HookRegistry) *Set {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	return &Set{
		Collection:   NewCollectionService(r),
		Field:        NewFieldService(r),
		Content:      NewContentService(r),
		ContentValue: NewContentValueService(r, hr),
		Asset:        NewAssetService(r),
		User:         NewUserService(r, []byte(secret)),
		Apikey:       NewApikeyService(r),
		Webhook:      NewWebhookService(r),
		Api:          NewApiService(r),
	}
}
