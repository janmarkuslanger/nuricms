package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
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
