package service

import (
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type ApiService interface {
	FindContentByCollectionAlias(alias string, offset int, perPage int) ([]dto.ContentItemResponse, error)
	FindContentByID(id uint) (dto.ContentItemResponse, error)
	FindContentByCollectionAndFieldValue(alias, fieldAlias, value string, offset, perPage int) ([]dto.ContentItemResponse, error)
	PrepareContent(ce *model.Content) (dto.ContentItemResponse, error)
}

type apiService struct {
	repos *repository.Set
}

func NewApiService(repos *repository.Set) ApiService {
	return &apiService{
		repos: repos,
	}
}

func (s *apiService) PrepareContent(ce *model.Content) (dto.ContentItemResponse, error) {
	values := make(map[string]any, len(ce.ContentValues))

	for _, cv := range ce.ContentValues {

		alias := cv.Field.Alias

		cvr := dto.ContentValueResponse{
			ID:        cv.ID,
			Value:     cv.Value,
			FieldType: cv.Field.FieldType,
		}

		if cv.Field.FieldType == model.FieldTypeCollection {
			id, _ := utils.StringToUint(cv.Value)
			con, err := s.repos.Content.FindByID(id)
			if err == nil {
				cvr.Collection = &dto.CollectionResponse{
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
				cvr.Asset = &dto.AssetResponse{
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

	return dto.ContentItemResponse{
		ID:        ce.ID,
		CreatedAt: ce.CreatedAt,
		UpdatedAt: ce.UpdatedAt,
		Values:    values,
		Collection: dto.CollectionResponse{
			Alias: ce.Collection.Alias,
			ID:    ce.CollectionID,
			Name:  ce.Collection.Name,
		},
	}, nil
}

func (s *apiService) FindContentByCollectionAlias(alias string, offset int, perPage int) ([]dto.ContentItemResponse, error) {
	var data []dto.ContentItemResponse

	collection, err := s.repos.Collection.FindByAlias(alias)
	if err != nil {
		return data, err
	}

	content, err := s.repos.Content.FindByCollectionID(collection.ID, offset, perPage)
	if err != nil {
		return data, err
	}

	for _, ce := range content {
		ci, err := s.PrepareContent(&ce)
		if err != nil {
			return data, nil
		}

		data = append(data, ci)
	}

	return data, nil
}

func (s *apiService) FindContentByID(id uint) (dto.ContentItemResponse, error) {
	var data dto.ContentItemResponse

	content, err := s.repos.Content.FindByID(id)
	if err != nil {
		return data, err
	}

	return s.PrepareContent(content)
}

func (s *apiService) FindContentByCollectionAndFieldValue(alias, fieldAlias, value string, offset, perPage int) ([]dto.ContentItemResponse, error) {
	collection, _ := s.repos.Collection.FindByAlias(alias)
	contents, _, _ := s.repos.Content.FindByCollectionAndFieldValue(collection.ID, fieldAlias, value, offset, perPage)
	var items []dto.ContentItemResponse
	for _, ce := range contents {
		ci, _ := s.PrepareContent(&ce)
		items = append(items, ci)
	}
	return items, nil
}
