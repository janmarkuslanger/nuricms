package service

import (
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

type ContentItemResponse struct {
	ID         uint               `json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Values     map[string]any     `json:"values"`
	Collection CollectionResponse `json:"collection"`
}

type ContentValueResponse struct {
	ID        uint            `json:"id"`
	Value     any             `json:"value"`
	FieldType model.FieldType `json:"field_type"`
}

// func (s *ApiService) transformContentRecursive(ce *model.Content) ContentItemResponse {
// 	contentValues := make(map[string]any)

// 	for _, cv := range ce.ContentValues {
// 		alias := cv.Field.Alias
// 		var val any

// 		switch cv.Field.FieldType {
// 		case model.FieldTypeCollection:
// 			id, _ := strconv.ParseUint(cv.Value, 10, 32)
// 			cont, _ := s.repos.Content.FindByID(uint(id))
// 			val = s.transformContentRecursive(&cont)

// 		case model.FieldTypeAsset:
// 			id, _ := strconv.ParseUint(cv.Value, 10, 32)
// 			asset, _ := s.repos.Asset.FindByID(uint(id))
// 			val = asset

// 		default:
// 			val = cv.Value
// 		}

// 		cvr := ContentValueResponse{
// 			ID:        cv.ID,
// 			Value:     val,
// 			FieldType: cv.Field.FieldType,
// 		}

// 		if cv.Field.IsList {
// 			slice, ok := contentValues[alias].([]any)
// 			if !ok {
// 				slice = []any{}
// 			}
// 			slice = append(slice, cvr)
// 			contentValues[alias] = slice

// 		} else {
// 			contentValues[alias] = cvr
// 		}
// 	}

// 	return ContentItemResponse{
// 		ID:        ce.ID,
// 		CreatedAt: ce.CreatedAt,
// 		UpdatedAt: ce.UpdatedAt,
// 		Values:    contentValues,
// 		Collection: CollectionResponse{
// 			ID:    ce.Collection.ID,
// 			Name:  ce.Collection.Name,
// 			Alias: ce.Collection.Alias,
// 		},
// 	}
// }

func (s *ApiService) FindContentByCollectionAlias(alias string, offset int, perPage int) ([]ContentItemResponse, error) {
	var data []ContentItemResponse

	collection, err := s.repos.Collection.FindByAlias(alias)
	if err != nil {
		return data, err
	}

	content, err := s.repos.Content.FindByCollectionID(collection.ID, offset, perPage)
	if err != nil {
		return data, err
	}

	for _, ce := range content {
		values := make(map[string]any, len(ce.ContentValues))

		for _, cv := range ce.ContentValues {

			alias := cv.Field.Alias

			cvr := ContentValueResponse{
				ID:        cv.ID,
				Value:     cv.Value,
				FieldType: cv.Field.FieldType,
			}

			if cv.Field.IsList {
				items, ok := values[alias].([]any)
				if !ok {
					items = []any{}
				}
				items = append(items, cvr)
				values[alias] = items

			} else {
				values[alias] = cvr
			}
		}

		data = append(data, ContentItemResponse{
			ID:        ce.ID,
			CreatedAt: ce.CreatedAt,
			UpdatedAt: ce.UpdatedAt,
			Values:    values,
		})
	}

	return data, nil
}
