package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
)

type mockCollectionRepo struct{ mock.Mock }

func (m *mockCollectionRepo) Create(c *model.Collection) error {
	return m.Called(c).Error(0)
}

func (m *mockCollectionRepo) Save(c *model.Collection) error {
	return m.Called(c).Error(0)
}

func (m *mockCollectionRepo) Delete(c *model.Collection) error {
	return m.Called(c).Error(0)
}

func (m *mockCollectionRepo) FindByID(id uint, opts ...base.QueryOption) (*model.Collection, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockCollectionRepo) List(page, pageSize int, opts ...base.QueryOption) ([]model.Collection, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *mockCollectionRepo) FindByAlias(alias string) (*model.Collection, error) {
	args := m.Called(alias)
	if obj := args.Get(0); obj != nil {
		return obj.(*model.Collection), args.Error(1)
	}
	return nil, args.Error(1)
}

func newTestCollectionService(repo repository.CollectionRepo) *CollectionService {
	return NewCollectionService(&repository.Set{Collection: repo})
}
func TestCollectionService_List(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := NewCollectionService(&repository.Set{Collection: repo})

	sample := []model.Collection{{Model: gorm.Model{ID: 1}}}
	repo.On("List", 2, 5).Return(sample, int64(1), nil)

	list, total, err := svc.List(2, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, sample, list)
}

func TestCollectionService_List_Error(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := NewCollectionService(&repository.Set{Collection: repo})

	repo.On("List", 1, 1).Return([]model.Collection{}, int64(0), errors.New("fail"))

	list, total, err := svc.List(1, 1)
	assert.EqualError(t, err, "fail")
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}

func TestCollectionService_FindByID(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	col := &model.Collection{Model: gorm.Model{ID: 7}}
	repo.On("FindByID", uint(7)).Return(col, nil)

	got, err := svc.FindByID(7)
	assert.NoError(t, err)
	assert.Equal(t, col, got)
}

func TestCollectionService_FindByID_NotFound(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	repo.On("FindByID", uint(9)).Return(nil, errors.New("nf"))

	_, err := svc.FindByID(9)
	assert.EqualError(t, err, "nf")
}

func TestCollectionService_FindByAlias(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	col := &model.Collection{Model: gorm.Model{ID: 3}, Alias: "a"}
	repo.On("FindByAlias", "a").Return(col, nil)

	got, err := svc.FindByAlias("a")
	assert.NoError(t, err)
	assert.Equal(t, col, got)
}

func TestCollectionService_FindByAlias_NotFound(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	repo.On("FindByAlias", "x").Return(nil, errors.New("nf2"))

	_, err := svc.FindByAlias("x")
	assert.EqualError(t, err, "nf2")
}

func TestCollectionService_Create(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	data := &dto.CollectionData{Name: "N", Alias: "A", Description: "D"}
	repo.On("Create", mock.MatchedBy(func(c *model.Collection) bool {
		return c.Name == "N" && c.Alias == "A" && c.Description == "D"
	})).Return(nil)

	c, err := svc.Create(data)
	assert.NoError(t, err)
	assert.Equal(t, "N", c.Name)
	assert.Equal(t, "A", c.Alias)
	assert.Equal(t, "D", c.Description)
}

func TestCollectionService_Create_Error(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	data := &dto.CollectionData{Name: "N", Alias: "A", Description: ""}
	repo.On("Create", mock.Anything).Return(errors.New("cfail"))

	_, err := svc.Create(data)
	assert.EqualError(t, err, "cfail")
}

func TestCollectionService_DeleteByID(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	col := &model.Collection{Model: gorm.Model{ID: 5}}
	repo.On("FindByID", uint(5)).Return(col, nil)
	repo.On("Delete", col).Return(nil)

	err := svc.DeleteByID(5)
	assert.NoError(t, err)
}

func TestCollectionService_DeleteByID_FindError(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	repo.On("FindByID", uint(6)).Return(nil, errors.New("dx"))

	err := svc.DeleteByID(6)
	assert.EqualError(t, err, "dx")
}

func TestCollectionService_UpdateByID_Success(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	col := &model.Collection{Model: gorm.Model{ID: 8}, Name: "Old", Alias: "OA", Description: "OD"}
	repo.On("FindByID", uint(8)).Return(col, nil)
	data := dto.CollectionData{Name: "New", Alias: "NA", Description: "ND"}
	repo.On("Save", col).Return(nil)

	updated, err := svc.UpdateByID(8, data)
	assert.NoError(t, err)
	assert.Equal(t, "New", updated.Name)
	assert.Equal(t, "NA", updated.Alias)
	assert.Equal(t, "ND", updated.Description)
}

func TestCollectionService_UpdateByID_FindError(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	repo.On("FindByID", uint(9)).Return(nil, errors.New("dx2"))

	_, err := svc.UpdateByID(9, dto.CollectionData{})
	assert.EqualError(t, err, "dx2")
}

func TestCollectionService_UpdateByID_NoAlias(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	col := &model.Collection{Model: gorm.Model{ID: 10}}
	repo.On("FindByID", uint(10)).Return(col, nil)

	_, err := svc.UpdateByID(10, dto.CollectionData{Name: "NameOnly"})
	assert.EqualError(t, err, "no alias given")
}

func TestCollectionService_UpdateByID_NoName(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	col := &model.Collection{Model: gorm.Model{ID: 11}, Alias: "AliasOnly"}
	repo.On("FindByID", uint(11)).Return(col, nil)

	_, err := svc.UpdateByID(11, dto.CollectionData{Alias: "AliasOnly"})
	assert.EqualError(t, err, "no name given")
}

func TestCollectionService_UpdateByID_SaveError(t *testing.T) {
	repo := new(mockCollectionRepo)
	svc := newTestCollectionService(repo)

	col := &model.Collection{Model: gorm.Model{ID: 12}, Name: "X", Alias: "Y"}
	repo.On("FindByID", uint(12)).Return(col, nil)
	data := dto.CollectionData{Name: "X2", Alias: "Y2"}
	repo.On("Save", mock.MatchedBy(func(c *model.Collection) bool {
		return c.Name == "X2" && c.Alias == "Y2"
	})).Return(errors.New("sfail"))

	_, err := svc.UpdateByID(12, data)
	assert.EqualError(t, err, "sfail")
}
