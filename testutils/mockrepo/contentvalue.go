package mockrepo

import (
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/janmarkuslanger/nuricms/pkg/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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

func (m *MockContentValueRepo) WithTx(tx *gorm.DB) repository.ContentValueRepo {
	m.Called(tx)
	return m
}
