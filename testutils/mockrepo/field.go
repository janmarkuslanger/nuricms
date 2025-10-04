package mockrepo

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockFieldRepo struct {
	mock.Mock
}

func (m *MockFieldRepo) Create(field *model.Field) error {
	args := m.Called(field)
	return args.Error(0)
}

func (m *MockFieldRepo) Save(field *model.Field) error {
	args := m.Called(field)
	return args.Error(0)
}

func (m *MockFieldRepo) Delete(field *model.Field) error {
	args := m.Called(field)
	return args.Error(0)
}

func (m *MockFieldRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Field, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Field), args.Error(1)
}

func (m *MockFieldRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Field, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Field), args.Get(1).(int64), args.Error(2)
}

func (m *MockFieldRepo) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFieldRepo) FindByCollectionID(collectionID uint) ([]model.Field, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.Field), args.Error(1)
}

func (m *MockFieldRepo) FindByFieldTypes(fieldTypes []model.FieldType) ([]model.Field, error) {
	args := m.Called(fieldTypes)
	return args.Get(0).([]model.Field), args.Error(1)
}

func (m *MockFieldRepo) FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.Field), args.Error(1)
}

func (m *MockFieldRepo) WithTx(tx *gorm.DB) repository.FieldRepo {
	m.Called(tx)
	return m
}
