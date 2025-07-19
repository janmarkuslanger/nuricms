package setup_test

import (
	"errors"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/setup"
	"github.com/janmarkuslanger/nuricms/testutils/mockservices"
	"github.com/stretchr/testify/assert"
)

func TestInitAdminUser_ListErr(t *testing.T) {
	s := &mockservices.MockUserService{}
	s.On("List", 1, 1).Return(make([]model.User, 1), int64(1), errors.New("something"))
	err := setup.InitAdminUser(s)
	assert.EqualError(t, err, "something")
}

func TestInitAdminUser_HasUser(t *testing.T) {
	s := &mockservices.MockUserService{}
	s.On("List", 1, 1).Return(make([]model.User, 1), int64(1), nil)
	err := setup.InitAdminUser(s)
	assert.Equal(t, err, nil)
}

func TestInitAdminUser_CreateErr(t *testing.T) {
	s := &mockservices.MockUserService{}
	s.On("List", 1, 1).Return(make([]model.User, 0), int64(0), nil)
	s.On("Create", dto.UserData{
		Email:    "admin@admin.com",
		Password: "mysecret",
		Role:     string(model.RoleAdmin),
	}).Return(&model.User{}, errors.New("create failed"))
	err := setup.InitAdminUser(s)
	assert.EqualError(t, err, "create failed")
}

func TestInitAdminUser_Success(t *testing.T) {
	s := &mockservices.MockUserService{}
	s.On("List", 1, 1).Return(make([]model.User, 0), int64(0), nil)
	s.On("Create", dto.UserData{
		Email:    "admin@admin.com",
		Password: "mysecret",
		Role:     string(model.RoleAdmin),
	}).Return(&model.User{}, nil)
	err := setup.InitAdminUser(s)
	assert.Equal(t, err, nil)
}
