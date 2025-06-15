package testutils

import (
	"time"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockApikeyService struct{ mock.Mock }

func (m *MockApikeyService) List(page, pageSize int) ([]model.Apikey, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Apikey), args.Get(1).(int64), args.Error(2)
}
func (m *MockApikeyService) Create(name string, ttl time.Duration) (string, error) {
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
