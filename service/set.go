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
	Apikey       *ApikeyService
}

func NewSet(r *repository.Set) *Set {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	return &Set{
		Collection:   NewCollectionService(r),
		Field:        NewFieldService(r),
		Content:      NewContentService(r),
		ContentValue: NewContentValueService(r),
		Asset:        NewAssetService(r),
		User:         NewUserService(r, []byte(secret)),
		Apikey:       NewApikeyService(r),
	}
}
