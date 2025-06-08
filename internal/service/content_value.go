package service

import (
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
)

type ContentValueService struct {
	repos        *repository.Set
	hookRegistry *plugin.HookRegistry
}

func NewContentValueService(repos *repository.Set, hr *plugin.HookRegistry) *ContentValueService {
	return &ContentValueService{repos: repos, hookRegistry: hr}
}

func (s *ContentValueService) Create(cv *model.ContentValue) error {
	s.hookRegistry.Run("contentValue:beforeSave", cv)
	return s.repos.ContentValue.Create(cv)
}
