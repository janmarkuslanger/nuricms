package service

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type FieldService interface {
	DeleteByID(id uint) error
	FindByCollectionID(collectionID uint) ([]model.Field, error)
	FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error)
	FindByID(id uint) (*model.Field, error)
	List(page, pageSize int) ([]model.Field, int64, error)
	Create(data dto.FieldData) (*model.Field, error)
	UpdateByID(id uint, data dto.FieldData) (*model.Field, error)
}

type fieldService struct {
	repos *repository.Set
}

func NewFieldService(repos *repository.Set) FieldService {
	return &fieldService{repos: repos}
}

func (s *fieldService) FindByCollectionID(collectionID uint) ([]model.Field, error) {
	return s.repos.Field.FindByCollectionID(collectionID)
}

func (s *fieldService) FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	return s.repos.Field.FindDisplayFieldsByCollectionID(collectionID)
}

func (s *fieldService) FindByID(id uint) (*model.Field, error) {
	return s.repos.Field.FindByID(id)
}

func (s *fieldService) List(page, pageSize int) ([]model.Field, int64, error) {
	return s.repos.Field.List(page, pageSize)
}

func (s *fieldService) UpdateByID(id uint, data dto.FieldData) (*model.Field, error) {
	field, err := s.repos.Field.FindByID(id)
	if err != nil {
		return nil, err
	}

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

	field.Name = data.Name
	field.Alias = data.Alias
	field.CollectionID = collectionID
	field.FieldType = model.FieldType(data.FieldType)
	field.IsList = data.IsList == "on"
	field.IsRequired = data.IsRequired == "on"
	field.DisplayField = data.DisplayField == "on"

	err = s.repos.Field.Save(field)
	return field, err
}

func (s *fieldService) Create(data dto.FieldData) (*model.Field, error) {
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

func (s *fieldService) DeleteByID(id uint) error {
	field, err := s.repos.Field.FindByID(id)
	if err != nil {
		return err
	}

	return s.repos.Field.Delete(field)
}
