package service

import (
	"errors"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/janmarkuslanger/nuricms/internal/repository/base"
)

type mockContentValueRepo struct {
	mock.Mock
}

func (m *mockContentValueRepo) Create(cv *model.ContentValue) error {
	return m.Called(cv).Error(0)
}
func (m *mockContentValueRepo) Save(cv *model.ContentValue) error {
	return m.Called(cv).Error(0)
}
func (m *mockContentValueRepo) Delete(cv *model.ContentValue) error {
	return m.Called(cv).Error(0)
}
func (m *mockContentValueRepo) FindByContentID(id uint) ([]model.ContentValue, error) {
	args := m.Called(id)
	return args.Get(0).([]model.ContentValue), args.Error(1)
}
func (m *mockContentValueRepo) FindByID(id uint, opts ...base.QueryOption) (*model.ContentValue, error) {
	return nil, nil
}
func (m *mockContentValueRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.ContentValue, int64, error) {
	return nil, 0, nil
}

func newTestContentValueService(repo *mockContentValueRepo, hr *plugin.HookRegistry) ContentValueService {
	return NewContentValueService(&repository.Set{ContentValue: repo}, hr)
}

func TestContentValueService_Create_RunsHookAndCreates(t *testing.T) {
	repo := new(mockContentValueRepo)
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
	repo := new(mockContentValueRepo)
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
	repo := new(mockContentValueRepo)
	hr := plugin.NewHookRegistry()
	hr.Register("contentValue:beforeSave", func(payload any) error { return nil })

	svc := newTestContentValueService(repo, hr)
	cv := &model.ContentValue{ContentID: 5, FieldID: 6, Value: "y"}
	repo.On("Create", cv).Return(errors.New("db error"))

	err := svc.Create(cv)
	assert.EqualError(t, err, "db error")
	assert.True(t, hr.Run("contentValue:beforeSave", cv) == nil)
}
