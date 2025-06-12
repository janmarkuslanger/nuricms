package service

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type FieldService struct {
	repos *repository.Set
}

func NewFieldService(repos *repository.Set) *FieldService {
	return &FieldService{repos: repos}
}

func (s *FieldService) GetByCollectionID(collectionID uint) ([]model.Field, error) {
	return s.repos.Field.FindByCollectionID(collectionID)
}

func (s *FieldService) GetDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	return s.repos.Field.FindDisplayFieldsByCollectionID(collectionID)
}

func (s *FieldService) List(page, pageSize int) ([]model.Field, int64, error) {
	return s.repos.Field.List(page, pageSize)
}

func (s *FieldService) Create(data dto.FieldData) (*model.Field, error) {
	collectionID, ok := utils.StringToUint(data.CollectionID)
	if !ok {
		return nil, errors.New("cannot convert collection id")
	}

	if data.Name == "" {
		return nil, errors.New("no name given")
	}

	if data.Alias == "" {
		return nil, errors.New("no alias given")
	}

	field := model.Field{
		Name:         data.Name,
		Alias:        data.Alias,
		CollectionID: collectionID,
		FieldType:    model.FieldType(data.FieldType),
		IsList:       data.IsList == "on",
		IsRequired:   data.IsRequired == "on",
		DisplayField: data.DisplayField == "on",
	}

	err := s.repos.Field.Create(&field)
	return &field, err
}
