package mockservices

import (
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockUserService) Create(dto dto.UserData) (*model.User, error) {
	args := m.Called(dto)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateByID(id uint, data dto.UserData) (*model.User, error) {
	args := m.Called(id, data)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
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
