package mockrepo

import (
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/stretchr/testify/mock"
)

type MockWebhookRepo struct {
	mock.Mock
}

func (m *MockWebhookRepo) Create(hook *model.Webhook) error {
	args := m.Called(hook)
	return args.Error(0)
}

func (m *MockWebhookRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Webhook, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Webhook), args.Error(1)
}

func (m *MockWebhookRepo) Update(hook *model.Webhook) error {
	args := m.Called(hook)
	return args.Error(0)
}

func (m *MockWebhookRepo) Delete(hook *model.Webhook) error {
	args := m.Called(hook)
	return args.Error(0)
}

func (m *MockWebhookRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Webhook, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Webhook), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebhookRepo) ListByEvent(event string) ([]model.Webhook, error) {
	args := m.Called(event)
	return args.Get(0).([]model.Webhook), args.Error(1)
}

func (m *MockWebhookRepo) Save(item *model.Webhook) error {
	args := m.Called(item)
	return args.Error(0)
}
