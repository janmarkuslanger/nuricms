package service

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/db"
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"gorm.io/gorm"
)

type ContentService interface {
	DeleteContentValuesByID(id uint) error
	EditWithValues(cwv dto.ContentWithValues) (*model.Content, error)
	DeleteByID(id uint) error
	CreateWithValues(cwv dto.ContentWithValues) (*model.Content, error)
	FindContentsWithDisplayContentValue() ([]model.Content, error)
	FindDisplayValueByCollectionID(collectionID uint, page, pageSize int) ([]model.Content, int64, error)
	FindByCollectionID(collectionID uint) ([]model.Content, error)
	ListByCollectionAlias(alias string, offset int, limit int) ([]model.Content, error)
	FindByID(id uint) (*model.Content, error)
	Create(c *model.Content) (*model.Content, error)
}

type contentService struct {
	repos *repository.Set
}

func NewContentService(repos *repository.Set) *contentService {
	return &contentService{repos: repos}
}

func (s *contentService) Create(c *model.Content) (*model.Content, error) {
	return c, s.repos.Content.Create(c)
}

func (s *contentService) FindByID(id uint) (*model.Content, error) {
	return s.repos.Content.FindByID(id)
}

func (s *contentService) ListByCollectionAlias(alias string, offset int, limit int) ([]model.Content, error) {
	var contents []model.Content

	collection, err := s.repos.Collection.FindByAlias(alias)
	if err != nil {
		return contents, err
	}

	return s.repos.Content.FindByCollectionID(collection.ID, offset, limit)
}

func (s *contentService) FindByCollectionID(collectionID uint) ([]model.Content, error) {
	return s.repos.Content.FindByCollectionID(collectionID, 0, 0)
}

func (s *contentService) FindDisplayValueByCollectionID(collectionID uint, page, pageSize int) ([]model.Content, int64, error) {
	return s.repos.Content.FindDisplayValueByCollectionID(collectionID, page, pageSize)
}

func (s *contentService) FindContentsWithDisplayContentValue() ([]model.Content, error) {
	return s.repos.Content.ListWithDisplayContentValue()
}

func (s *contentService) CreateWithValues(cwv dto.ContentWithValues) (*model.Content, error) {
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

func (s *contentService) DeleteByID(id uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		err := s.repos.Content.DeleteByID(id)
		if err != nil {
			return err
		}

		err = s.DeleteContentValuesByID(id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *contentService) DeleteContentValuesByID(id uint) error {
	values, err := s.repos.ContentValue.FindByContentID(id)
	if err != nil {
		return err
	}

	for _, v := range values {
		err = s.repos.ContentValue.Delete(&v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *contentService) EditWithValues(cwv dto.ContentWithValues) (*model.Content, error) {
	var content model.Content
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		content, err := s.repos.Content.FindByID(cwv.ContentID)
		if err != nil {
			return err
		}

		if content.CollectionID != cwv.CollectionID {
			return errors.New("content doesnt relate to Collection")
		}

		if err = s.DeleteContentValuesByID(content.ID); err != nil {
			return err
		}

		fields, err := s.repos.Field.FindByCollectionID(cwv.CollectionID)
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
