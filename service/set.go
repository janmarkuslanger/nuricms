package service

import "github.com/janmarkuslanger/nuricms/repository"

type Set struct {
	Collection   *CollectionService
	Field        *FieldService
	Content      *ContentService
	ContentValue *ContentValueService
	Asset        *AssetService
}

func NewSet(r *repository.Set) *Set {
	return &Set{
		Collection:   NewCollectionService(r.Collection),
		Field:        NewFieldService(r.Field),
		Content:      NewContentService(r.Content),
		ContentValue: NewContentValueService(r.ContentValue),
		Asset:        NewAssetService(r.Asset),
	}
}
