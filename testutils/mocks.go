package testutils

import (
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/fs"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/stretchr/testify/mock"
)

type MockApikeyRepo struct {
	mock.Mock
}

func (m *MockApikeyRepo) Create(a *model.Apikey) error {
	return m.Called(a).Error(0)
}

func (m *MockApikeyRepo) Delete(a *model.Apikey) error {
	return m.Called(a).Error(0)
}

func (m *MockApikeyRepo) Save(a *model.Apikey) error {
	return m.Called(a).Error(0)
}

func (m *MockApikeyRepo) FindByToken(token string) (*model.Apikey, error) {
	args := m.Called(token)
	if val := args.Get(0); val != nil {
		return val.(*model.Apikey), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockApikeyRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Apikey, error) {
	args := m.Called(id)
	if val := args.Get(0); val != nil {
		return val.(*model.Apikey), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockApikeyRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Apikey, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Apikey), args.Get(1).(int64), args.Error(2)
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
func (m *MockCollectionService) Create(data dto.CollectionData) (*model.Collection, error) {
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
func (m *MockCollectionService) Save(col *model.Collection) error {
	args := m.Called(col)
	return args.Error(0)
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

func (m *MockAssetService) UploadFile(ctx server.Context, file fs.FileOpener, filename string) (string, error) {
	args := m.Called(ctx, file, filename)
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

type MockWebhookService struct {
	mock.Mock
}

func (m *MockWebhookService) Create(dto dto.WebhookData) (*model.Webhook, error) {
	args := m.Called(dto)
	return args.Get(0).(*model.Webhook), args.Error(1)
}

func (m *MockWebhookService) List(page, pageSize int) ([]model.Webhook, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Webhook), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebhookService) UpdateByID(id uint, data dto.WebhookData) (*model.Webhook, error) {
	args := m.Called(id, data)
	if webhook := args.Get(0); webhook != nil {
		return webhook.(*model.Webhook), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWebhookService) FindByID(id uint) (*model.Webhook, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Webhook), args.Error(1)
}

func (m *MockWebhookService) Save(webhook *model.Webhook) error {
	args := m.Called(webhook)
	return args.Error(0)
}

func (m *MockWebhookService) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWebhookService) Dispatch(event string, payload any) error {
	m.Called(event, payload)
	return nil
}

type MockContentService struct {
	mock.Mock
}

func (m *MockContentService) EditWithValues(cwv dto.ContentWithValues) (*model.Content, error) {
	args := m.Called(cwv)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Content), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContentService) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockContentService) CreateWithValues(cwv dto.ContentWithValues) (*model.Content, error) {
	args := m.Called(cwv)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Content), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContentService) FindContentsWithDisplayContentValue() ([]model.Content, error) {
	args := m.Called()
	return args.Get(0).([]model.Content), args.Error(1)
}

func (m *MockContentService) FindDisplayValueByCollectionID(collectionID uint, page, pageSize int) ([]model.Content, int64, error) {
	args := m.Called(collectionID, page, pageSize)
	return args.Get(0).([]model.Content), args.Get(1).(int64), args.Error(2)
}

func (m *MockContentService) FindByCollectionID(collectionID uint) ([]model.Content, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.Content), args.Error(1)
}

func (m *MockContentService) ListByCollectionAlias(alias string, offset int, limit int) ([]model.Content, error) {
	args := m.Called(alias, offset, limit)
	return args.Get(0).([]model.Content), args.Error(1)
}

func (m *MockContentService) FindByID(id uint) (*model.Content, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Content), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContentService) Create(c *model.Content) (*model.Content, error) {
	args := m.Called(c)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Content), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockApiService struct {
	mock.Mock
}

func (m *MockApiService) FindContentByCollectionAlias(alias string, offset int, perPage int) ([]dto.ContentItemResponse, error) {
	args := m.Called(alias, offset, perPage)
	return args.Get(0).([]dto.ContentItemResponse), args.Error(1)
}

func (m *MockApiService) FindContentByID(id uint) (dto.ContentItemResponse, error) {
	args := m.Called(id)
	return args.Get(0).(dto.ContentItemResponse), args.Error(1)
}

func (m *MockApiService) FindContentByCollectionAndFieldValue(alias, fieldAlias, value string, offset, perPage int) ([]dto.ContentItemResponse, error) {
	args := m.Called(alias, fieldAlias, value, offset, perPage)
	return args.Get(0).([]dto.ContentItemResponse), args.Error(1)
}

func (m *MockApiService) PrepareContent(content *model.Content) (dto.ContentItemResponse, error) {
	args := m.Called(content)
	return args.Get(0).(dto.ContentItemResponse), args.Error(1)
}

type MockApikeyService struct {
	mock.Mock
}

func (m *MockApikeyService) List(page, pageSize int) ([]model.Apikey, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Apikey), args.Get(1).(int64), args.Error(2)
}

func (m *MockApikeyService) Create(data dto.ApikeyData) (*model.Apikey, error) {
	args := m.Called(data)
	return args.Get(0).(*model.Apikey), args.Error(1)
}

func (m *MockApikeyService) FindByID(id uint) (*model.Apikey, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Apikey), args.Error(1)
}

func (m *MockApikeyService) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockApikeyService) Validate(token string) error {
	args := m.Called(token)
	return args.Error(0)
}
