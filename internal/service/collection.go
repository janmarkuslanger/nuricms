package service

import (
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
)

type CollectionService struct {
	repos *repository.Set
}

func NewCollectionService(repos *repository.Set) *CollectionService {
	return &CollectionService{repos: repos}
}

func (s *CollectionService) List(page, pageSize int) ([]model.Collection, int64, error) {
	return s.repos.Collection.List(page, pageSize)
}

func (s *CollectionService) FindByID(id uint) (*model.Collection, error) {
	return s.repos.Collection.FindByID(id)
}

func (s *CollectionService) FindByAlias(alias string) (*model.Collection, error) {
	return s.repos.Collection.FindByAlias(alias)
}

func (s *CollectionService) Create(data *dto.CollectionData) (*model.Collection, error) {
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
