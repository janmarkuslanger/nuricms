package service

import (
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/model"
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
