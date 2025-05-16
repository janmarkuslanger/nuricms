package service

import (
	"log"
	"os"

	"github.com/janmarkuslanger/nuricms/repository"
)

type Set struct {
	Collection   *CollectionService
	Field        *FieldService
	Content      *ContentService
	ContentValue *ContentValueService
	Asset        *AssetService
	User         *UserService
}

func NewSet(r *repository.Set) *Set {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	return &Set{
		Collection:   NewCollectionService(r.Collection),
		Field:        NewFieldService(r.Field),
		Content:      NewContentService(r.Content),
		ContentValue: NewContentValueService(r.ContentValue),
		Asset:        NewAssetService(r.Asset),
		User:         NewUserService(r.User, []byte(secret)),
	}
}
