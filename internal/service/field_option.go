package service

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
)

type FieldOptionService interface {
	Create(fieldOption model.FieldOption) (*model.FieldOption, error)
}

type fieldOptionService struct {
	repos *repository.Set
}

func NewFieldOptionService(repos *repository.Set) *fieldOptionService {
	return &fieldOptionService{repos: repos}
}
