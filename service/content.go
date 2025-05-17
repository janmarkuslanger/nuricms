package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type ContentService struct {
	repos *repository.Set
}

func NewContentService(repos *repository.Set) *ContentService {
	return &ContentService{repos: repos}
}

func (s *ContentService) Create(c *model.Content) (model.Content, error) {
	return s.repos.Content.Create(c)
}

func (s *ContentService) FindByID(id uint) (model.Content, error) {
	return s.repos.Content.FindByID(id)
}

func (s *ContentService) ListByCollectionAlias(alias string, offset int, limit int) ([]model.Content, error) {
	var contents []model.Content

	collection, err := s.repos.Collection.FindByAlias(alias)
	if err != nil {
		return contents, err
	}

	return s.repos.Content.FindByCollectionID(collection.ID, offset, limit)
}

func (s *ContentService) GetByCollectionID(collectionID uint) ([]model.Content, error) {
	return s.repos.Content.FindByCollectionID(collectionID, 0, 0)
}

func (s *ContentService) GetDisplayValueByCollectionID(collectionID uint) ([]model.Content, error) {
	return s.repos.Content.FindDisplayValueByCollectionID(collectionID)
}

func (s *ContentService) GetContentsWithDisplayContentValue() ([]model.Content, error) {
	return s.repos.Content.ListWithDisplayContentValue()
}
