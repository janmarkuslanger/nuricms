package service

import (
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type ContentValueService struct {
	repo *repository.ContentValueRepository
}

func NewContentValueService(r *repository.ContentValueRepository) *ContentValueService {
	return &ContentValueService{repo: r}
}

func (s *ContentValueService) Create(cv *model.ContentValue) error {
	return s.repo.Create(cv)
}
