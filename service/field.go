package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type FieldService struct {
	repo *repository.FieldRepository
}

func NewFieldService(r *repository.FieldRepository) *FieldService {
	return &FieldService{repo: r}
}

func (s *FieldService) GetByCollectionID(collectionID uint) ([]model.Field, error) {
	return s.repo.FindByCollectionID(collectionID)
}

func (s *FieldService) GetDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	return s.repo.FindDisplayFieldsByCollectionID(collectionID)
}
