package testutils

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/stretchr/testify/mock"
)

type MockFieldRepo struct {
	mock.Mock
}

func (m *MockFieldRepo) Create(entity *model.Field) error {
	return m.Called(entity).Error(0)
}

func (m *MockFieldRepo) Save(entity *model.Field) error {
	return m.Called(entity).Error(0)
}

func (m *MockFieldRepo) Delete(entity *model.Field) error {
	return m.Called(entity).Error(0)
}

func (m *MockFieldRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Field, error) {
	args := m.Called(id)
	if val := args.Get(0); val != nil {
		return val.(*model.Field), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFieldRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Field, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Field), args.Get(1).(int64), args.Error(2)
}

func (m *MockFieldRepo) FindByCollectionID(collectionID uint) ([]model.Field, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.Field), args.Error(1)
}

func (m *MockFieldRepo) FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.Field), args.Error(1)
}

type MockCollectionRepo struct {
	mock.Mock
}

func (m *MockCollectionRepo) Create(c *model.Collection) error {
	return m.Called(c).Error(0)
}

func (m *MockCollectionRepo) Save(c *model.Collection) error {
	return m.Called(c).Error(0)
}

func (m *MockCollectionRepo) Delete(c *model.Collection) error {
	return m.Called(c).Error(0)
}

func (m *MockCollectionRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Collection, error) {
	args := m.Called(id)
	if val := args.Get(0); val != nil {
		return val.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCollectionRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Collection, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepo) FindByAlias(alias string) (*model.Collection, error) {
	args := m.Called(alias)
	if val := args.Get(0); val != nil {
		return val.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}
