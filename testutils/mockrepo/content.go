package mockrepo

import (
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/janmarkuslanger/nuricms/pkg/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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

func (m *MockContentRepo) WithTx(tx *gorm.DB) repository.ContentRepo {
	m.Called(tx)
	return m
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
