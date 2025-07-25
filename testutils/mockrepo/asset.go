package mockrepo

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/stretchr/testify/mock"
)

type MockAssetRepo struct {
	mock.Mock
}

func (m *MockAssetRepo) Create(entity *model.Asset) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockAssetRepo) Save(entity *model.Asset) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockAssetRepo) Delete(entity *model.Asset) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockAssetRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Asset, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Asset), args.Error(1)
}

func (m *MockAssetRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Asset, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Asset), args.Get(1).(int64), args.Error(2)
}
