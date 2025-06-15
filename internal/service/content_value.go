package service

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
)

type ContentValueService interface {
	Create(cv *model.ContentValue) error
}

type contentValueService struct {
	repos        *repository.Set
	hookRegistry *plugin.HookRegistry
}

func NewContentValueService(repos *repository.Set, hr *plugin.HookRegistry) ContentValueService {
	return &contentValueService{repos: repos, hookRegistry: hr}
}

func (s *contentValueService) Create(cv *model.ContentValue) error {
	s.hookRegistry.Run("contentValue:beforeSave", cv)
	return s.repos.ContentValue.Create(cv)
}
