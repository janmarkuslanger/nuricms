package service

import (
	"errors"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"github.com/janmarkuslanger/nuricms/testutils/mockrepo"
	"github.com/stretchr/testify/assert"
)

func newTestContentValueService(repo *mockrepo.MockContentValueRepo, hr *plugin.HookRegistry) ContentValueService {
	return NewContentValueService(&repository.Set{ContentValue: repo}, hr)
}

func TestContentValueService_Create_RunsHookAndCreates(t *testing.T) {
	repo := new(mockrepo.MockContentValueRepo)
	hr := plugin.NewHookRegistry()
	var invoked bool
	var gotPayload any
	hr.Register("contentValue:beforeSave", func(payload any) error {
		invoked = true
		gotPayload = payload
		return nil
	})

	svc := newTestContentValueService(repo, hr)
	cv := &model.ContentValue{ContentID: 1, FieldID: 2, Value: "v"}
	repo.On("Create", cv).Return(nil)

	err := svc.Create(cv)
	assert.NoError(t, err)
	assert.True(t, invoked)
	assert.Equal(t, cv, gotPayload)
	repo.AssertCalled(t, "Create", cv)
}

func TestContentValueService_Create_HookErrorIgnored(t *testing.T) {
	repo := new(mockrepo.MockContentValueRepo)
	hr := plugin.NewHookRegistry()
	hr.Register("contentValue:beforeSave", func(payload any) error {
		return errors.New("hookfail")
	})

	svc := newTestContentValueService(repo, hr)
	cv := &model.ContentValue{ContentID: 3, FieldID: 4, Value: "x"}
	repo.On("Create", cv).Return(nil)

	err := svc.Create(cv)
	assert.NoError(t, err)
	repo.AssertCalled(t, "Create", cv)
}

func TestContentValueService_Create_RepoError(t *testing.T) {
	repo := new(mockrepo.MockContentValueRepo)
	hr := plugin.NewHookRegistry()
	hr.Register("contentValue:beforeSave", func(payload any) error { return nil })

	svc := newTestContentValueService(repo, hr)
	cv := &model.ContentValue{ContentID: 5, FieldID: 6, Value: "y"}
	repo.On("Create", cv).Return(errors.New("db error"))

	err := svc.Create(cv)
	assert.EqualError(t, err, "db error")
	assert.True(t, hr.Run("contentValue:beforeSave", cv) == nil)
}
