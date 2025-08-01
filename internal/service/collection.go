package service

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
)

type CollectionService interface {
	UpdateByID(colID uint, data dto.CollectionData) (*model.Collection, error)
	DeleteByID(colOD uint) error
	Create(data dto.CollectionData) (*model.Collection, error)
	FindByAlias(alias string) (*model.Collection, error)
	FindByID(id uint) (*model.Collection, error)
	List(page, pageSize int) ([]model.Collection, int64, error)
	Save(col *model.Collection) error
}

type collectionService struct {
	repos *repository.Set
}

func NewCollectionService(repos *repository.Set) CollectionService {
	return &collectionService{repos: repos}
}

func (s *collectionService) List(page, pageSize int) ([]model.Collection, int64, error) {
	return s.repos.Collection.List(page, pageSize)
}

func (s *collectionService) FindByID(id uint) (*model.Collection, error) {
	return s.repos.Collection.FindByID(id)
}

func (s *collectionService) Save(col *model.Collection) error {
	return s.repos.Collection.Save(col)
}

func (s *collectionService) FindByAlias(alias string) (*model.Collection, error) {
	return s.repos.Collection.FindByAlias(alias)
}

func (s *collectionService) Create(data dto.CollectionData) (*model.Collection, error) {
	collection := &model.Collection{
		Name:        data.Name,
		Alias:       data.Alias,
		Description: data.Description,
	}

	err := s.repos.Collection.Create(collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (s *collectionService) DeleteByID(collectionID uint) error {
	collection, err := s.FindByID(uint(collectionID))
	if err != nil {
		return err
	}

	return s.repos.Collection.Delete(collection)
}

func (s *collectionService) UpdateByID(colID uint, data dto.CollectionData) (*model.Collection, error) {
	collection, err := s.FindByID(colID)
	if err != nil {
		return nil, err
	}

	if data.Alias == "" {
		return nil, errors.New("no alias given")
	}

	if data.Name == "" {
		return nil, errors.New("no name given")
	}

	collection.Alias = data.Alias
	collection.Name = data.Name
	collection.Description = data.Description

	err = s.repos.Collection.Save(collection)
	return collection, err
}
