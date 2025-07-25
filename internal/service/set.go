// service/set.go

package service

import (
	"github.com/janmarkuslanger/nuricms/internal/env"
	"github.com/janmarkuslanger/nuricms/internal/fs"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"gorm.io/gorm"
)

type Set struct {
	Collection   CollectionService
	Field        FieldService
	Content      ContentService
	ContentValue ContentValueService
	Asset        AssetService
	User         UserService
	Apikey       ApikeyService
	Webhook      WebhookService
	Api          ApiService
}

func NewSet(r *repository.Set, hr *plugin.HookRegistry, db *gorm.DB, env *env.Env, fs fs.FileOps) (*Set, error) {
	return &Set{
		Collection:   NewCollectionService(r),
		Field:        NewFieldService(r),
		Content:      NewContentService(r, db),
		ContentValue: NewContentValueService(r, hr),
		Asset:        NewAssetService(r, fs),
		User:         NewUserService(r, []byte(env.Secret)),
		Apikey:       NewApikeyService(r),
		Webhook:      NewWebhookService(r),
		Api:          NewApiService(r),
	}, nil
}
