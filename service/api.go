package service

import (
	"time"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
	"github.com/janmarkuslanger/nuricms/utils"
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
	ID    uint   `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Alias string `json:"alias,omitempty"`
}

type AssetResponse struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type ContentItemResponse struct {
	ID         uint               `json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Values     map[string]any     `json:"values"`
	Collection CollectionResponse `json:"collection"`
}

type ContentValueResponse struct {
	ID         uint                `json:"id"`
	Value      any                 `json:"value"`
	FieldType  model.FieldType     `json:"field_type"`
	Collection *CollectionResponse `json:"collection,omitempty"`
	Asset      *AssetResponse      `json:"asset,omitempty"`
}

func (s *ApiService) prepareContent(ce *model.Content) (ContentItemResponse, error) {
	values := make(map[string]any, len(ce.ContentValues))

	for _, cv := range ce.ContentValues {

		alias := cv.Field.Alias

		cvr := ContentValueResponse{
			ID:        cv.ID,
			Value:     cv.Value,
			FieldType: cv.Field.FieldType,
		}

		if cv.Field.FieldType == model.FieldTypeCollection {
			id, _ := utils.StringToUint(cv.Value)
			con, err := s.repos.Content.FindByID(id)
			if err == nil {
				cvr.Collection = &CollectionResponse{
					ID:    con.CollectionID,
					Name:  con.Collection.Name,
					Alias: con.Collection.Alias,
				}
			}
		}

		if cv.Field.FieldType == model.FieldTypeAsset {
			id, _ := utils.StringToUint(cv.Value)
			ass, err := s.repos.Asset.FindByID(id)
			if err == nil {
				cvr.Asset = &AssetResponse{
					ID:   ass.ID,
					Name: ass.Name,
					Path: ass.Path,
				}
			}
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

	return ContentItemResponse{
		ID:        ce.ID,
		CreatedAt: ce.CreatedAt,
		UpdatedAt: ce.UpdatedAt,
		Values:    values,
		Collection: CollectionResponse{
			Alias: ce.Collection.Alias,
			ID:    ce.CollectionID,
			Name:  ce.Collection.Name,
		},
	}, nil
}

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
		ci, err := s.prepareContent(&ce)
		if err != nil {
			return data, nil
		}

		data = append(data, ci)
	}

	return data, nil
}

func (s *ApiService) FindContentByID(id uint) (ContentItemResponse, error) {
	var data ContentItemResponse

	content, err := s.repos.Content.FindByID(id)
	if err != nil {
		return data, err
	}

	return s.prepareContent(&content)
}

func (s *ApiService) FindContentByCollectionAndFieldValue(alias, fieldAlias, value string, offset, perPage int) ([]ContentItemResponse, error) {
	collection, _ := s.repos.Collection.FindByAlias(alias)
	contents, _, _ := s.repos.Content.FindByCollectionAndFieldValue(collection.ID, fieldAlias, value, offset, perPage)
	var items []ContentItemResponse
	for _, ce := range contents {
		ci, _ := s.prepareContent(&ce)
		items = append(items, ci)
	}
	return items, nil
}
