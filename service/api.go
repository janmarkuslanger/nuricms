package service

import (
	"strconv"
	"time"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type ApiService struct {
	repos *repository.Set
}

func NewApiService(repos *repository.Set) *ApiService {
	return &ApiService{
		repos: repos,
	}
}

type CollectionResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type ContentResponse struct {
	Collection CollectionResponse    `json:"collection"`
	Items      []ContentItemResponse `json:"items"`
}

type ContentItemResponse struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Values    map[string]any `json:"values"`
}

type ContentValueResponse struct {
	ID    uint `json:"id"`
	Value any  `json:"value"`
}

func (s *ApiService) listByCollectionAlias(alias string, offset int, limit int) ([]model.Content, error) {
	var contents []model.Content

	collection, err := s.repos.Collection.FindByAlias(alias)
	if err != nil {
		return contents, err
	}

	return s.repos.Content.FindByCollectionID(collection.ID, offset, limit)
}

func (s *ApiService) transformContentRecursive(ce *model.Content) ContentItemResponse {
	contentValues := make(map[string]any)

	for _, cv := range ce.ContentValues {
		alias := cv.Field.Alias
		var val any

		switch cv.Field.FieldType {
		case model.FieldTypeCollection:
			id, _ := strconv.ParseUint(cv.Value, 10, 32)
			cont, _ := s.repos.Content.FindByID(uint(id))
			val = s.transformContentRecursive(&cont)

		case model.FieldTypeAsset:
			id, _ := strconv.ParseUint(cv.Value, 10, 32)
			asset, _ := s.repos.Asset.FindByID(uint(id))
			val = asset

		default:
			val = cv.Value
		}

		if cv.Field.IsList {
			slice, ok := contentValues[alias].([]any)
			if !ok {
				slice = []any{}
			}
			slice = append(slice, val)
			contentValues[alias] = slice

		} else {
			contentValues[alias] = val
		}
	}

	return ContentItemResponse{
		ID:        ce.ID,
		CreatedAt: ce.CreatedAt,
		UpdatedAt: ce.UpdatedAt,
		Values:    contentValues,
	}
}

func (s *ApiService) ListByCollectionAlias(alias string, offset int, perPage int) (ContentResponse, error) {
	var out ContentResponse

	col, err := s.repos.Collection.FindByAlias(alias)

	if err != nil {
		return out, err
	}

	content, err := s.listByCollectionAlias(alias, offset, perPage)
	if err != nil {
		return out, err
	}

	var contentItems []ContentItemResponse

	for _, ce := range content {
		contentItems = append(contentItems, s.transformContentRecursive(&ce))
	}

	out = ContentResponse{
		Collection: CollectionResponse{
			ID:    col.ID,
			Name:  col.Name,
			Alias: col.Alias,
		},
		Items: contentItems,
	}

	return out, nil
}
