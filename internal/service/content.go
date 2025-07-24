package service

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"gorm.io/gorm"
)

type ContentService interface {
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
	db    *gorm.DB
}

func NewContentService(repos *repository.Set, db *gorm.DB) *contentService {
	return &contentService{repos: repos, db: db}
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

func (s *contentService) saveContentValues(contentValueRepo repository.ContentValueRepo, contentID uint, fields []model.Field, formData map[string][]string) error {
	for _, f := range fields {
		for i, v := range formData[f.Alias] {
			cv := model.ContentValue{
				SortIndex: i + 1,
				ContentID: contentID,
				FieldID:   f.ID,
				Value:     v,
			}

			if err := contentValueRepo.Create(&cv); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *contentService) CreateWithValues(cwv dto.ContentWithValues) (*model.Content, error) {
	var content model.Content
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txField := s.repos.Field.WithTx(tx)
		txContent := s.repos.Content.WithTx(tx)
		txContentValue := s.repos.ContentValue.WithTx(tx)

		fields, err := txField.FindByCollectionID(cwv.CollectionID)
		if err != nil {
			return err
		}

		found := model.Content{CollectionID: cwv.CollectionID}
		content = found

		err = txContent.Create(&content)
		if err != nil {
			return err
		}

		if err := s.saveContentValues(txContentValue, content.ID, fields, cwv.FormData); err != nil {
			return err
		}

		return nil
	})

	return &content, err
}

func (s *contentService) DeleteByID(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txContent := s.repos.Content.WithTx(tx)
		txContentValue := s.repos.ContentValue.WithTx(tx)
		err := txContent.DeleteByID(id)
		if err != nil {
			return err
		}

		err = s.deleteContentValuesByID(txContentValue, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *contentService) deleteContentValuesByID(repo repository.ContentValueRepo, id uint) error {
	values, err := repo.FindByContentID(id)
	if err != nil {
		return err
	}

	for _, v := range values {
		err = repo.Delete(&v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *contentService) EditWithValues(cwv dto.ContentWithValues) (*model.Content, error) {
	var content *model.Content
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txContent := s.repos.Content.WithTx(tx)
		txContentValue := s.repos.ContentValue.WithTx(tx)

		found, err := txContent.FindByID(cwv.ContentID)
		content = found
		if err != nil {
			return err
		}

		if content.CollectionID != cwv.CollectionID {
			return errors.New("content doesnt relate to Collection")
		}

		if err = s.deleteContentValuesByID(txContentValue, content.ID); err != nil {
			return err
		}

		fields, err := s.repos.Field.FindByCollectionID(cwv.CollectionID)
		if err != nil {
			return err
		}

		if err := s.saveContentValues(txContentValue, content.ID, fields, cwv.FormData); err != nil {
			return err
		}

		return nil
	})

	return content, err
}
