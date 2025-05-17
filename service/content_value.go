package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type ContentValueService struct {
	repos *repository.Set
}

func NewContentValueService(repos *repository.Set) *ContentValueService {
	return &ContentValueService{repos: repos}
}

func (s *ContentValueService) Create(cv *model.ContentValue) error {
	return s.repos.ContentValue.Create(cv)
}
