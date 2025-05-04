package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type ContentService struct {
	repo *repository.ContentRepository
}

func NewContentService(r *repository.ContentRepository) *ContentService {
	return &ContentService{repo: r}
}

func (s *ContentService) Create(c *model.Content) (model.Content, error) {
	return s.repo.Create(c)
}

func (s *ContentService) GetByID(id uint) (model.Content, error) {
	return s.repo.FindByID(id)
}

func (s *ContentService) GetByCollectionID(collectionID uint) ([]model.Content, error) {
	return s.repo.FindByCollectionID(collectionID)
}

func (s *ContentService) GetDisplayValueByCollectionID(collectionID uint) ([]model.Content, error) {
	return s.repo.FindDisplayValueByCollectionID(collectionID)
}

func (s *ContentService) GetContentsWithDisplayContentValue() ([]model.Content, error) {
	return s.repo.FindAllWithDisplayContentValue()
}
