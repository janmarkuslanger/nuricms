package testutils

import (
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockApikeyService struct{ mock.Mock }

func (m *MockApikeyService) List(page, pageSize int) ([]model.Apikey, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Apikey), args.Get(1).(int64), args.Error(2)
}
func (m *MockApikeyService) CreateToken(name string, ttl time.Duration) (string, error) {
	args := m.Called(name, ttl)
	return args.String(0), args.Error(1)
}
func (m *MockApikeyService) FindByID(id uint) (*model.Apikey, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Apikey), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockApikeyService) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockApikeyService) Validate(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

type MockCollectionService struct{ mock.Mock }

func (m *MockCollectionService) UpdateByID(collectionID uint, data dto.CollectionData) (*model.Collection, error) {
	args := m.Called(collectionID, data)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCollectionService) DeleteByID(collectionID uint) error {
	args := m.Called(collectionID)
	return args.Error(0)
}
func (m *MockCollectionService) Create(data *dto.CollectionData) (*model.Collection, error) {
	args := m.Called(data)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCollectionService) FindByAlias(alias string) (*model.Collection, error) {
	args := m.Called(alias)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCollectionService) FindByID(id uint) (*model.Collection, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCollectionService) List(page, pageSize int) ([]model.Collection, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Collection), args.Get(1).(int64), args.Error(2)
}

type MockAssetService struct{ mock.Mock }

func (m *MockAssetService) List(page, pageSize int) ([]model.Asset, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Asset), args.Get(1).(int64), args.Error(2)
}

func (m *MockAssetService) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAssetService) Save(asset *model.Asset) error {
	args := m.Called(asset)
	return args.Error(0)
}

func (m *MockAssetService) Create(asset *model.Asset) error {
	args := m.Called(asset)
	return args.Error(0)
}

func (m *MockAssetService) UploadFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	args := m.Called(c, file)
	return args.String(0), args.Error(1)
}

func (m *MockAssetService) FindByID(id uint) (*model.Asset, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockFieldService struct{ mock.Mock }

func (m *MockFieldService) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFieldService) FindByCollectionID(collectionID uint) ([]model.Field, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.Field), args.Error(1)
}

func (m *MockFieldService) FindDisplayFieldsByCollectionID(collectionID uint) ([]model.Field, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.Field), args.Error(1)
}

func (m *MockFieldService) FindByID(id uint) (*model.Field, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Field), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFieldService) List(page, pageSize int) ([]model.Field, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Field), args.Get(1).(int64), args.Error(2)
}

func (m *MockFieldService) Create(data dto.FieldData) (*model.Field, error) {
	args := m.Called(data)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Field), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFieldService) UpdateByID(id uint, data dto.FieldData) (*model.Field, error) {
	args := m.Called(id, data)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Field), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) List(page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) Create(email, password string, role model.Role) (*model.User, error) {
	args := m.Called(email, password, role)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Save(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Delete(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) LoginUser(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) ValidateJWT(tokenStr string) (uint, string, model.Role, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(uint), args.String(1), args.Get(2).(model.Role), args.Error(3)
}
