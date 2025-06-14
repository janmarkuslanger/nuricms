package service

import (
	"github.com/janmarkuslanger/nuricms/internal/db"
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"gorm.io/gorm"
)

type ContentService struct {
	repos *repository.Set
}

func NewContentService(repos *repository.Set) *ContentService {
	return &ContentService{repos: repos}
}

func (s *ContentService) Create(c *model.Content) (*model.Content, error) {
	return c, s.repos.Content.Create(c)
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

func (s *ContentService) FindByCollectionID(collectionID uint) ([]model.Content, error) {
	return s.repos.Content.FindByCollectionID(collectionID, 0, 0)
}

func (s *ContentService) FindDisplayValueByCollectionID(collectionID uint, page, pageSize int) ([]model.Content, int64, error) {
	return s.repos.Content.FindDisplayValueByCollectionID(collectionID, page, pageSize)
}

func (s *ContentService) FindContentsWithDisplayContentValue() ([]model.Content, error) {
	return s.repos.Content.ListWithDisplayContentValue()
}

func (s *ContentService) CreateWithValues(cwv dto.ContentWithValues) (*model.Content, error) {
	var content model.Content
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		fields, err := s.repos.Field.FindByCollectionID(cwv.CollectionID)
		if err != nil {
			return err
		}

		content = model.Content{CollectionID: cwv.CollectionID}

		err = s.repos.Content.Create(&content)
		if err != nil {
			return err
		}

		for _, f := range fields {
			for i, v := range cwv.FormData[f.Alias] {
				cv := model.ContentValue{
					SortIndex: i + 1,
					ContentID: content.ID,
					FieldID:   f.ID,
					Value:     v,
				}

				if err := s.repos.ContentValue.Create(&cv); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return &content, err
}
