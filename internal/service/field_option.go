package service

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/utils"
	"gorm.io/gorm"
)

type FieldOptionService interface {
	Create(dto dto.FieldOption) (*model.FieldOption, error)
	List(page, pageSize int) ([]model.FieldOption, int64, error)
	UpdateByID(id uint, dto dto.FieldOption) (*model.FieldOption, error)
	FindByID(id uint) (*model.FieldOption, error)
	DeleteByID(id uint) error
}

type fieldOptionService struct {
	repos *repository.Set
}

func NewFieldOptionService(repos *repository.Set) *fieldOptionService {
	return &fieldOptionService{repos: repos}
}

func (s *fieldOptionService) FindByID(id uint) (*model.FieldOption, error) {
	return s.repos.FieldOption.FindByID(id, func(db *gorm.DB) *gorm.DB {
		return db.Preload("Field")
	})
}

func (s *fieldOptionService) DeleteByID(id uint) error {
	item, err := s.repos.FieldOption.FindByID(id)
	if err != nil {
		return err
	}

	return s.repos.FieldOption.Delete(item)
}

func (s *fieldOptionService) Create(dto dto.FieldOption) (*model.FieldOption, error) {
	fieldID, ok := utils.StringToUint(dto.FieldID)
	if !ok {
		return nil, errors.New("cannot convert field id")
	}

	if dto.Value == "" {
		return nil, errors.New("no value given")
	}

	fo := model.FieldOption{
		Value:      dto.Value,
		FieldID:    fieldID,
		OptionType: model.FieldOptionType(dto.OptionType),
	}

	err := s.repos.FieldOption.Create(&fo)
	return &fo, err
}

func (s *fieldOptionService) UpdateByID(id uint, dto dto.FieldOption) (*model.FieldOption, error) {
	item, err := s.repos.FieldOption.FindByID(id)
	if err != nil {
		return nil, errors.New("no field found")
	}

	// we dont allow changing the field
	item.Value = dto.Value

	err = s.repos.FieldOption.Save(item)
	if err != nil {
		return nil, errors.New("could not save field option")
	}

	return item, nil
}

func (s *fieldOptionService) List(page, pageSize int) ([]model.FieldOption, int64, error) {
	return s.repos.FieldOption.List(page, pageSize, func(db *gorm.DB) *gorm.DB {
		return db.Preload("Field").Preload("Field.Collection")
	})
}
