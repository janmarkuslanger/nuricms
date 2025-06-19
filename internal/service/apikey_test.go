package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
)

type mockApikeyRepo struct{ mock.Mock }

func (m *mockApikeyRepo) Create(a *model.Apikey) error {
	return m.Called(a).Error(0)
}
func (m *mockApikeyRepo) Save(a *model.Apikey) error {
	return m.Called(a).Error(0)
}
func (m *mockApikeyRepo) Delete(a *model.Apikey) error {
	return m.Called(a).Error(0)
}
func (m *mockApikeyRepo) FindByToken(token string) (*model.Apikey, error) {
	args := m.Called(token)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Apikey), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockApikeyRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Apikey, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Apikey), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockApikeyRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Apikey, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Apikey), args.Get(1).(int64), args.Error(2)
}

func newTestService(repo repository.ApikeyRepo) ApikeyService {
	set := &repository.Set{Apikey: repo}
	return NewApikeyService(set)
}

func TestApikeyService_Create_Success(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	repo.
		On("Create", mock.MatchedBy(func(a *model.Apikey) bool {
			return a.Name == "test" && len(a.Token) == 64
		})).
		Return(nil)
	token, err := svc.CreateToken("test", 0)
	assert.NoError(t, err)
	assert.Len(t, token, 64)
	repo.AssertExpectations(t)
}

func TestApikeyService_Create_FailureOnRepo(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	repo.On("Create", mock.Anything).Return(errors.New("db error"))
	token, err := svc.CreateToken("xyz", time.Minute)
	assert.Empty(t, token)
	assert.EqualError(t, err, "db error")
}

func TestApikeyService_Validate_Success(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	future := time.Now().Add(time.Hour)
	key := &model.Apikey{Model: gorm.Model{ID: 1}, Token: "tok", ExpiresAt: &future}
	repo.On("FindByToken", "tok").Return(key, nil)
	err := svc.Validate("tok")
	assert.NoError(t, err)
}

func TestApikeyService_Validate_InvalidToken(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	repo.On("FindByToken", "bad").Return(nil, errors.New("not found"))
	err := svc.Validate("bad")
	assert.EqualError(t, err, "invalid api key")
}

func TestApikeyService_Validate_Expired(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	past := time.Now().Add(-time.Hour)
	key := &model.Apikey{Model: gorm.Model{ID: 2}, Token: "tok2", ExpiresAt: &past}
	repo.On("FindByToken", "tok2").Return(key, nil)
	err := svc.Validate("tok2")
	assert.EqualError(t, err, "api key expired")
}

func TestApikeyService_ListAndFindByID(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	sample := []model.Apikey{
		{Model: gorm.Model{ID: 1}},
		{Model: gorm.Model{ID: 2}},
	}
	repo.On("List", 1, 2).Return(sample, int64(2), nil)
	list, total, err := svc.List(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, sample, list)

	key := &model.Apikey{Model: gorm.Model{ID: 3}}
	repo.On("FindByID", uint(3)).Return(key, nil)
	found, err := svc.FindByID(3)
	assert.NoError(t, err)
	assert.Equal(t, key, found)
}

func TestApikeyService_DeleteByID_Success(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	key := &model.Apikey{Model: gorm.Model{ID: 5}}
	repo.On("FindByID", uint(5)).Return(key, nil)
	repo.On("Delete", key).Return(nil)
	err := svc.DeleteByID(5)
	assert.NoError(t, err)
}

func TestApikeyService_DeleteByID_NotFound(t *testing.T) {
	repo := new(mockApikeyRepo)
	svc := newTestService(repo)
	repo.On("FindByID", uint(7)).Return(nil, errors.New("not found"))
	err := svc.DeleteByID(7)
	assert.EqualError(t, err, "not found")
}
