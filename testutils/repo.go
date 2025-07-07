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

type MockContentRepo struct {
	mock.Mock
}

func (m *MockContentRepo) Create(entity *model.Content) error {
	return m.Called(entity).Error(0)
}

func (m *MockContentRepo) Save(entity *model.Content) error {
	return m.Called(entity).Error(0)
}

func (m *MockContentRepo) Delete(entity *model.Content) error {
	return m.Called(entity).Error(0)
}

func (m *MockContentRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Content, error) {
	args := m.Called(id)
	if val := args.Get(0); val != nil {
		return val.(*model.Content), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContentRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Content, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Content), args.Get(1).(int64), args.Error(2)
}

func (m *MockContentRepo) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockContentRepo) FindByCollectionID(collectionID uint, offset, limit int) ([]model.Content, error) {
	args := m.Called(collectionID, offset, limit)
	return args.Get(0).([]model.Content), args.Error(1)
}

func (m *MockContentRepo) FindDisplayValueByCollectionID(collectionID uint, page, pageSize int) ([]model.Content, int64, error) {
	args := m.Called(collectionID, page, pageSize)
	return args.Get(0).([]model.Content), args.Get(1).(int64), args.Error(2)
}

func (m *MockContentRepo) ListWithDisplayContentValue() ([]model.Content, error) {
	args := m.Called()
	return args.Get(0).([]model.Content), args.Error(1)
}

func (m *MockContentRepo) FindByCollectionAndFieldValue(collectionID uint, fieldAlias, value string, offset, limit int) ([]model.Content, int, error) {
	args := m.Called(collectionID, fieldAlias, value, offset, limit)
	return args.Get(0).([]model.Content), args.Int(1), args.Error(2)
}

type MockContentValueRepo struct {
	mock.Mock
}

func (m *MockContentValueRepo) Create(entity *model.ContentValue) error {
	return m.Called(entity).Error(0)
}

func (m *MockContentValueRepo) Save(entity *model.ContentValue) error {
	return m.Called(entity).Error(0)
}

func (m *MockContentValueRepo) Delete(entity *model.ContentValue) error {
	return m.Called(entity).Error(0)
}

func (m *MockContentValueRepo) FindByID(id uint, opts ...base.QueryOption) (*model.ContentValue, error) {
	args := m.Called(id)
	if val := args.Get(0); val != nil {
		return val.(*model.ContentValue), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContentValueRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.ContentValue, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.ContentValue), args.Get(1).(int64), args.Error(2)
}

func (m *MockContentValueRepo) FindByContentID(cID uint) ([]model.ContentValue, error) {
	args := m.Called(cID)
	return args.Get(0).([]model.ContentValue), args.Error(1)
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
