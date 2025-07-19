package service_test

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteByID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(testutils.MockContentRepo)
	mockContentValueRepo := new(testutils.MockContentValueRepo)

	repos := &repository.Set{
		Content:      mockContentRepo,
		ContentValue: mockContentValueRepo,
	}

	s := service.NewContentService(repos, testDB)

	id := uint(1)
	mockContentRepo.On("DeleteByID", id).Return(nil)
	mockContentValueRepo.On("FindByContentID", id).Return([]model.ContentValue{{}}, nil)
	mockContentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil)

	err := s.DeleteByID(id)
	assert.NoError(t, err)
}

func TestDeleteContentValuesByID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentValueRepo := new(testutils.MockContentValueRepo)
	repos := &repository.Set{ContentValue: mockContentValueRepo}
	s := service.NewContentService(repos, testDB)

	mockContentValueRepo.On("FindByContentID", uint(1)).Return([]model.ContentValue{{}}, nil)
	mockContentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil)

	err := s.DeleteContentValuesByID(1)
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(testutils.MockContentRepo)
	repos := &repository.Set{Content: mockContentRepo}
	s := service.NewContentService(repos, testDB)

	input := &model.Content{CollectionID: 1}
	mockContentRepo.On("Create", input).Return(nil)

	out, err := s.Create(input)
	assert.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestFindByID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(testutils.MockContentRepo)
	repos := &repository.Set{Content: mockContentRepo}
	s := service.NewContentService(repos, testDB)

	mockContent := &model.Content{}
	mockContentRepo.On("FindByID", uint(42)).Return(mockContent, nil)

	result, err := s.FindByID(42)
	assert.NoError(t, err)
	assert.Equal(t, mockContent, result)
}

func TestListByCollectionAlias(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(testutils.MockContentRepo)
	mockCollectionRepo := new(testutils.MockCollectionRepo)
	repos := &repository.Set{
		Content:    mockContentRepo,
		Collection: mockCollectionRepo,
	}
	s := service.NewContentService(repos, testDB)

	collection := &model.Collection{}
	collection.ID = 1
	mockCollectionRepo.On("FindByAlias", "blog").Return(collection, nil)
	mockContentRepo.On("FindByCollectionID", uint(1), 0, 10).Return([]model.Content{{}}, nil)

	contents, err := s.ListByCollectionAlias("blog", 0, 10)
	assert.NoError(t, err)
	assert.Len(t, contents, 1)
}

func TestFindByCollectionID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(testutils.MockContentRepo)
	repos := &repository.Set{Content: mockContentRepo}
	s := service.NewContentService(repos, testDB)

	mockContentRepo.On("FindByCollectionID", uint(1), 0, 0).Return([]model.Content{{}}, nil)
	result, err := s.FindByCollectionID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestFindDisplayValueByCollectionID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(testutils.MockContentRepo)
	repos := &repository.Set{Content: mockContentRepo}
	s := service.NewContentService(repos, testDB)

	mockContentRepo.On("FindDisplayValueByCollectionID", uint(1), 0, 10).Return([]model.Content{{}}, int64(1), nil)
	result, count, err := s.FindDisplayValueByCollectionID(1, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), count)
}

func TestFindContentsWithDisplayContentValue(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(testutils.MockContentRepo)
	repos := &repository.Set{Content: mockContentRepo}
	s := service.NewContentService(repos, testDB)

	mockContentRepo.On("ListWithDisplayContentValue").Return([]model.Content{{}}, nil)
	result, err := s.FindContentsWithDisplayContentValue()
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
