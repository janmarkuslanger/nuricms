package service_test

import (
	"errors"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/janmarkuslanger/nuricms/testutils/mockrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeleteByID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
	mockContentValueRepo := new(mockrepo.MockContentValueRepo)

	repos := &repository.Set{
		Content:      mockContentRepo,
		ContentValue: mockContentValueRepo,
	}

	s := service.NewContentService(repos, testDB)

	id := uint(1)
	mockContentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentRepo)
	mockContentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentValueRepo)
	mockContentRepo.On("DeleteByID", id).Return(nil)
	mockContentValueRepo.On("FindByContentID", id).Return([]model.ContentValue{{}}, nil)
	mockContentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil)

	err := s.DeleteByID(id)
	assert.NoError(t, err)
}

func TestDeleteByID_Failed(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
	mockContentValueRepo := new(mockrepo.MockContentValueRepo)

	repos := &repository.Set{
		Content:      mockContentRepo,
		ContentValue: mockContentValueRepo,
	}

	s := service.NewContentService(repos, testDB)

	id := uint(1)
	mockContentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentRepo)
	mockContentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentValueRepo)
	mockContentRepo.On("DeleteByID", id).Return(errors.New("error"))
	mockContentValueRepo.On("FindByContentID", id).Return([]model.ContentValue{{}}, nil)
	mockContentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil)

	err := s.DeleteByID(id)
	assert.Error(t, err)
}

func TestDeleteByID_FindByContentIDErr(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
	mockContentValueRepo := new(mockrepo.MockContentValueRepo)

	repos := &repository.Set{
		Content:      mockContentRepo,
		ContentValue: mockContentValueRepo,
	}

	s := service.NewContentService(repos, testDB)

	id := uint(1)
	mockContentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentRepo)
	mockContentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentValueRepo)
	mockContentRepo.On("DeleteByID", id).Return(nil)
	mockContentValueRepo.On("FindByContentID", id).Return([]model.ContentValue{{}}, errors.New("error"))
	mockContentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil)

	err := s.DeleteByID(id)
	assert.Error(t, err)
}

func TestDeleteByID_DeleteErr(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
	mockContentValueRepo := new(mockrepo.MockContentValueRepo)

	repos := &repository.Set{
		Content:      mockContentRepo,
		ContentValue: mockContentValueRepo,
	}

	s := service.NewContentService(repos, testDB)

	id := uint(1)
	mockContentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentRepo)
	mockContentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockContentValueRepo)
	mockContentRepo.On("DeleteByID", id).Return(nil)
	mockContentValueRepo.On("FindByContentID", id).Return([]model.ContentValue{{}}, nil)
	mockContentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(errors.New("failed"))

	err := s.DeleteByID(id)
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
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
	mockContentRepo := new(mockrepo.MockContentRepo)
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
	mockContentRepo := new(mockrepo.MockContentRepo)
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

func TestListByCollectionAlias_FindByAliasErr(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
	mockCollectionRepo := new(testutils.MockCollectionRepo)
	repos := &repository.Set{
		Content:    mockContentRepo,
		Collection: mockCollectionRepo,
	}
	s := service.NewContentService(repos, testDB)

	collection := &model.Collection{}
	collection.ID = 1
	mockCollectionRepo.On("FindByAlias", "blog").Return(collection, errors.New("Err"))

	contents, err := s.ListByCollectionAlias("blog", 0, 10)
	assert.NotNil(t, err)
	assert.Len(t, contents, 0)
}

func TestFindByCollectionID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
	repos := &repository.Set{Content: mockContentRepo}
	s := service.NewContentService(repos, testDB)

	mockContentRepo.On("FindByCollectionID", uint(1), 0, 0).Return([]model.Content{{}}, nil)
	result, err := s.FindByCollectionID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestFindDisplayValueByCollectionID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	mockContentRepo := new(mockrepo.MockContentRepo)
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
	mockContentRepo := new(mockrepo.MockContentRepo)
	repos := &repository.Set{Content: mockContentRepo}
	s := service.NewContentService(repos, testDB)

	mockContentRepo.On("ListWithDisplayContentValue").Return([]model.Content{{}}, nil)
	result, err := s.FindContentsWithDisplayContentValue()
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func setupWithValues(t *testing.T) (*gorm.DB, *mockrepo.MockFieldRepo, *mockrepo.MockContentRepo, *mockrepo.MockContentValueRepo, *repository.Set) {
	testDB := testutils.SetupTestDB(t)

	fieldRepo := new(mockrepo.MockFieldRepo)
	contentRepo := new(mockrepo.MockContentRepo)
	contentValueRepo := new(mockrepo.MockContentValueRepo)

	repos := &repository.Set{
		Content:      contentRepo,
		Field:        fieldRepo,
		ContentValue: contentValueRepo,
	}

	return testDB, fieldRepo, contentRepo, contentValueRepo, repos
}

func TestCreateWithValues(t *testing.T) {
	testDB, fieldRepo, contentRepo, contentValueRepo, repos := setupWithValues(t)
	s := service.NewContentService(repos, testDB)

	fields := []model.Field{
		{
			Model: gorm.Model{ID: 1},
			Alias: "title",
		},
		{
			Model: gorm.Model{ID: 2},
			Alias: "desc",
		},
	}

	form := map[string][]string{
		"title": {"My Title"},
		"desc":  {"My Description"},
	}

	fieldRepo.On("WithTx", mock.Anything).Return(fieldRepo)
	fieldRepo.On("FindByCollectionID", uint(123)).Return(fields, nil)

	contentRepo.On("WithTx", mock.Anything).Return(contentRepo)
	contentRepo.On("Create", mock.AnythingOfType("*model.Content")).Return(nil)

	contentValueRepo.On("WithTx", mock.Anything).Return(contentValueRepo)
	contentValueRepo.On("Create", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()

	result, err := s.CreateWithValues(dto.ContentWithValues{
		CollectionID: 123,
		FormData:     form,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(123), result.CollectionID)

	fieldRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	contentValueRepo.AssertExpectations(t)
}

func TestCreateWithValues_FindByCollectionErr(t *testing.T) {
	testDB, fieldRepo, contentRepo, contentValueRepo, repos := setupWithValues(t)

	s := service.NewContentService(repos, testDB)

	fields := []model.Field{
		{
			Model: gorm.Model{ID: 1},
			Alias: "title",
		},
		{
			Model: gorm.Model{ID: 2},
			Alias: "desc",
		},
	}

	form := map[string][]string{
		"title": {"My Title"},
		"desc":  {"My Description"},
	}

	fieldRepo.On("WithTx", mock.Anything).Return(fieldRepo)
	fieldRepo.On("FindByCollectionID", uint(123)).Return(fields, errors.New("Something"))
	contentRepo.On("WithTx", mock.Anything).Return(contentRepo)
	contentValueRepo.On("WithTx", mock.Anything).Return(contentValueRepo)

	_, err := s.CreateWithValues(dto.ContentWithValues{
		CollectionID: 123,
		FormData:     form,
	})

	assert.NotNil(t, err)
	fieldRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	contentValueRepo.AssertExpectations(t)
}

func TestCreateWithValues_CreateErr(t *testing.T) {
	testDB, fieldRepo, contentRepo, contentValueRepo, repos := setupWithValues(t)
	s := service.NewContentService(repos, testDB)

	fields := []model.Field{
		{
			Model: gorm.Model{ID: 1},
			Alias: "title",
		},
		{
			Model: gorm.Model{ID: 2},
			Alias: "desc",
		},
	}

	form := map[string][]string{
		"title": {"My Title"},
		"desc":  {"My Description"},
	}

	fieldRepo.On("WithTx", mock.Anything).Return(fieldRepo)
	fieldRepo.On("FindByCollectionID", uint(123)).Return(fields, nil)

	contentRepo.On("WithTx", mock.Anything).Return(contentRepo)
	contentRepo.On("Create", mock.AnythingOfType("*model.Content")).Return(errors.New("ohno"))

	contentValueRepo.On("WithTx", mock.Anything).Return(contentValueRepo)
	contentValueRepo.On("Create", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()

	_, err := s.CreateWithValues(dto.ContentWithValues{
		CollectionID: 123,
		FormData:     form,
	})

	assert.NotNil(t, err)

	fieldRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	contentValueRepo.AssertExpectations(t)
}

func TestEditWithValues(t *testing.T) {
	testDB, fieldRepo, contentRepo, contentValueRepo, repos := setupWithValues(t)
	s := service.NewContentService(repos, testDB)

	fields := []model.Field{
		{Model: gorm.Model{ID: 1}, Alias: "title"},
		{Model: gorm.Model{ID: 2}, Alias: "desc"},
	}

	form := map[string][]string{
		"title": {"Updated Title"},
		"desc":  {"Updated Description"},
	}

	content := &model.Content{
		Model:        gorm.Model{ID: 42},
		CollectionID: 123,
	}

	contentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentRepo)
	contentRepo.On("FindByID", mock.Anything).Return(content, nil)

	contentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentValueRepo).Maybe()
	contentValueRepo.On("FindByContentID", uint(42)).Return([]model.ContentValue{}, nil)
	contentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()
	contentValueRepo.On("Create", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()

	fieldRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(fieldRepo).Maybe()
	fieldRepo.On("FindByCollectionID", uint(123)).Return(fields, nil)

	result, err := s.EditWithValues(dto.ContentWithValues{
		ContentID:    42,
		CollectionID: 123,
		FormData:     form,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, uint(42), result.ID)
		assert.Equal(t, uint(123), result.CollectionID)
	}

	fieldRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	contentValueRepo.AssertExpectations(t)
}

func TestEditWithValues_NotFound(t *testing.T) {
	testDB, fieldRepo, contentRepo, contentValueRepo, repos := setupWithValues(t)
	s := service.NewContentService(repos, testDB)
	form := map[string][]string{
		"title": {"Updated Title"},
		"desc":  {"Updated Description"},
	}

	content := &model.Content{
		Model:        gorm.Model{ID: 42},
		CollectionID: 123,
	}

	contentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentRepo)
	contentRepo.On("FindByID", mock.Anything).Return(content, errors.New("oh no"))

	contentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentValueRepo).Maybe()
	contentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()
	contentValueRepo.On("Create", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()

	fieldRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(fieldRepo).Maybe()

	_, err := s.EditWithValues(dto.ContentWithValues{
		ContentID:    42,
		CollectionID: 123,
		FormData:     form,
	})

	assert.Error(t, err)

	fieldRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	contentValueRepo.AssertExpectations(t)
}

func TestEditWithValues_CollectionIDInvalid(t *testing.T) {
	testDB, fieldRepo, contentRepo, contentValueRepo, repos := setupWithValues(t)
	s := service.NewContentService(repos, testDB)
	form := map[string][]string{
		"title": {"Updated Title"},
		"desc":  {"Updated Description"},
	}

	content := &model.Content{
		Model:        gorm.Model{ID: 42},
		CollectionID: 123,
	}

	contentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentRepo)
	contentRepo.On("FindByID", mock.Anything).Return(content, nil)

	contentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentValueRepo).Maybe()
	contentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()
	contentValueRepo.On("Create", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()

	fieldRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(fieldRepo).Maybe()

	_, err := s.EditWithValues(dto.ContentWithValues{
		ContentID:    42,
		CollectionID: 123 + 1,
		FormData:     form,
	})

	assert.Error(t, err)

	fieldRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	contentValueRepo.AssertExpectations(t)
}

func TestEditWithValues_NoFields(t *testing.T) {
	testDB, fieldRepo, contentRepo, contentValueRepo, repos := setupWithValues(t)
	s := service.NewContentService(repos, testDB)

	fields := []model.Field{
		{Model: gorm.Model{ID: 1}, Alias: "title"},
		{Model: gorm.Model{ID: 2}, Alias: "desc"},
	}

	form := map[string][]string{
		"title": {"Updated Title"},
		"desc":  {"Updated Description"},
	}

	content := &model.Content{
		Model:        gorm.Model{ID: 42},
		CollectionID: 123,
	}

	contentRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentRepo)
	contentRepo.On("FindByID", mock.Anything).Return(content, nil)

	contentValueRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(contentValueRepo).Maybe()
	contentValueRepo.On("FindByContentID", uint(42)).Return([]model.ContentValue{}, nil)
	contentValueRepo.On("Delete", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()
	contentValueRepo.On("Create", mock.AnythingOfType("*model.ContentValue")).Return(nil).Maybe()

	fieldRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(fieldRepo).Maybe()
	fieldRepo.On("FindByCollectionID", uint(123)).Return(fields, errors.New("no fields"))

	_, err := s.EditWithValues(dto.ContentWithValues{
		ContentID:    42,
		CollectionID: 123,
		FormData:     form,
	})

	assert.Error(t, err)

	fieldRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	contentValueRepo.AssertExpectations(t)
}
