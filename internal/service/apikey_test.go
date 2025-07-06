package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/testutils"
)

func newTestService(repo repository.ApikeyRepo) ApikeyService {
	return NewApikeyService(&repository.Set{Apikey: repo})
}

func TestCreate_Success(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	repo.On("Create", mock.MatchedBy(func(a *model.Apikey) bool {
		return a.Name == "Test" && len(a.Token) == 64
	})).Return(nil)

	key, err := svc.Create(dto.ApikeyData{Name: "Test"})

	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, "Test", key.Name)
	assert.Len(t, key.Token, 64)

	repo.AssertExpectations(t)
}

func TestCreate_NoName(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	key, err := svc.Create(dto.ApikeyData{Name: ""})

	assert.Nil(t, key)
	assert.EqualError(t, err, "no name given")
}

func TestCreate_RepoFails(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	repo.On("Create", mock.Anything).Return(errors.New("db error"))

	key, err := svc.Create(dto.ApikeyData{Name: "Fail"})

	assert.Nil(t, key)
	assert.EqualError(t, err, "db error")
}

func TestValidate_Success(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	key := &model.Apikey{
		Model: gorm.Model{ID: 1},
		Token: "abc",
	}
	repo.On("FindByToken", "abc").Return(key, nil)

	err := svc.Validate("abc")
	assert.NoError(t, err)
}

func TestValidate_Expired(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	past := time.Now().Add(-1 * time.Hour)
	key := &model.Apikey{
		Model:     gorm.Model{ID: 2},
		Token:     "expired",
		ExpiresAt: &past,
	}
	repo.On("FindByToken", "expired").Return(key, nil)

	err := svc.Validate("expired")
	assert.EqualError(t, err, "api key expired")
}

func TestValidate_NotFound(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	repo.On("FindByToken", "invalid").Return(nil, errors.New("not found"))

	err := svc.Validate("invalid")
	assert.EqualError(t, err, "invalid api key")
}

func TestList(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	mockList := []model.Apikey{
		{Model: gorm.Model{ID: 1}, Name: "One"},
		{Model: gorm.Model{ID: 2}, Name: "Two"},
	}
	repo.On("List", 1, 5).Return(mockList, int64(2), nil)

	list, total, err := svc.List(1, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, mockList, list)
}

func TestFindByID(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	key := &model.Apikey{Model: gorm.Model{ID: 10}, Name: "Test"}
	repo.On("FindByID", uint(10)).Return(key, nil)

	found, err := svc.FindByID(10)
	assert.NoError(t, err)
	assert.Equal(t, key, found)
}

func TestDeleteByID_Success(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	key := &model.Apikey{Model: gorm.Model{ID: 9}, Name: "DeleteMe"}
	repo.On("FindByID", uint(9)).Return(key, nil)
	repo.On("Delete", key).Return(nil)

	err := svc.DeleteByID(9)
	assert.NoError(t, err)
}

func TestDeleteByID_NotFound(t *testing.T) {
	repo := new(testutils.MockApikeyRepo)
	svc := newTestService(repo)

	repo.On("FindByID", uint(404)).Return(nil, errors.New("not found"))

	err := svc.DeleteByID(404)
	assert.EqualError(t, err, "not found")
}
