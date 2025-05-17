package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type CollectionService struct {
	repos *repository.Set
}

func NewCollectionService(repos *repository.Set) *CollectionService {
	return &CollectionService{repos: repos}
}

func (s *CollectionService) GetAll() ([]model.Collection, error) {
	return s.repos.Collection.GetAll()
}

func (s *CollectionService) GetByID(id uint) (*model.Collection, error) {
	return s.repos.Collection.FindByID(id)
}
