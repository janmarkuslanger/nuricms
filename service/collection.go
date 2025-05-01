package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type CollectionService struct {
	repo *repository.CollectionRepository
}

func NewCollectionService(r *repository.CollectionRepository) *CollectionService {
	return &CollectionService{repo: r}
}

func (s *CollectionService) GetAll() ([]model.Collection, error) {
	return s.repo.GetAll()
}

func (s *CollectionService) GetByID(id uint) (*model.Collection, error) {
	return s.repo.FindByID(id)
}
