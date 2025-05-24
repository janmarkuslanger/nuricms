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

func (s *CollectionService) List(page, pageSize int) ([]model.Collection, int64, error) {
	return s.repos.Collection.List(page, pageSize)
}

func (s *CollectionService) GetByID(id uint) (*model.Collection, error) {
	return s.repos.Collection.FindByID(id)
}
