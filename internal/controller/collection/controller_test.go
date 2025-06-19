package collection

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

type mockCollectionService struct {
	testutils.MockCollectionService
}

func createMockController(svc service.CollectionService) *Controller {
	gin.SetMode(gin.TestMode)
	svcSet := &service.Set{Collection: svc}
	return NewController(svcSet)
}

func TestShowCollections_Success(t *testing.T) {
	mockSvc := new(mockCollectionService)
	cols := []model.Collection{
		{Name: "A", Alias: "a"},
		{Name: "B", Alias: "b"},
	}
	mockSvc.On("List", 2, 5).Return(cols, int64(2), nil)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/collections?page=2&pageSize=5")
	ct.showCollections(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/index.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestShowCollections_InvalidQueryParams(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("List", 1, 10).Return([]model.Collection{}, int64(0), nil)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/collections?page=abc&pageSize=xyz")
	ct.showCollections(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/index.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestShowCollections_ListError(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("List", 1, 10).Return([]model.Collection{}, int64(0), errors.New("fail"))
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/collections")
	ct.showCollections(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/index.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestShowCreateCollection(t *testing.T) {
	mockSvc := new(mockCollectionService)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/collections/create")
	ct.showCreateCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
}

func TestCreateCollection_MissingFields(t *testing.T) {
	mockSvc := new(mockCollectionService)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakePOSTContext("/collections/create", gin.H{
		"name": "", "alias": "", "description": "Desc",
	})
	ct.createCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
	mockSvc.AssertNotCalled(t, "Create", mockSvc)
}

func TestCreateCollection_ServiceError(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("Create", &dto.CollectionData{
		Name:        "Name",
		Alias:       "alias",
		Description: "Desc",
	}).Return((*model.Collection)(nil), errors.New("fail"))
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakePOSTContext("/collections/create", gin.H{
		"name":        "Name",
		"alias":       "alias",
		"description": "Desc",
	})
	ct.createCollection(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestCreateCollection_Success(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("Create", &dto.CollectionData{
		Name:        "Name",
		Alias:       "alias",
		Description: "Desc",
	}).Return(&model.Collection{Name: "Name", Alias: "alias", Description: "Desc"}, nil)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/collections/create", gin.H{
		"name":        "Name",
		"alias":       "alias",
		"description": "Desc",
	})
	ct.createCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/collections", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestShowEditCollection_InvalidParam(t *testing.T) {
	mockSvc := new(mockCollectionService)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/collections/edit/abc")
	testutils.SetParam(c, "id", "abc")
	ct.showEditCollection(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
}

func TestShowEditCollection_NotFound(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("FindByID", uint(5)).Return((*model.Collection)(nil), errors.New("not found"))
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/collections/edit/5")
	testutils.SetParam(c, "id", "5")
	ct.showEditCollection(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestShowEditCollection_Success(t *testing.T) {
	mockSvc := new(mockCollectionService)
	col := &model.Collection{Name: "Name", Alias: "alias", Description: "Desc"}
	mockSvc.On("FindByID", uint(2)).Return(col, nil)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/collections/edit/2")
	testutils.SetParam(c, "id", "2")
	ct.showEditCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestEditCollection_InvalidParam(t *testing.T) {
	mockSvc := new(mockCollectionService)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakePOSTContext("/collections/edit/abc", gin.H{
		"name":        "Name",
		"alias":       "Alias",
		"description": "Desc",
	})
	testutils.SetParam(c, "id", "abc")
	ct.editCollection(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
}

func TestEditCollection_UpdateError(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("UpdateByID", uint(2), dto.CollectionData{
		Name:        "Name",
		Alias:       "Alias",
		Description: "Desc",
	}).Return((*model.Collection)(nil), errors.New("fail"))
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakePOSTContext("/collections/edit/2", gin.H{
		"name":        "Name",
		"alias":       "Alias",
		"description": "Desc",
	})
	testutils.SetParam(c, "id", "2")
	ct.editCollection(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:collection/create_or_edit.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestEditCollection_Success(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("UpdateByID", uint(3), dto.CollectionData{
		Name:        "Name",
		Alias:       "Alias",
		Description: "Desc",
	}).Return(&model.Collection{Name: "Name", Alias: "Alias", Description: "Desc"}, nil)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/collections/edit/3", gin.H{
		"name":        "Name",
		"alias":       "Alias",
		"description": "Desc",
	})
	testutils.SetParam(c, "id", "3")
	ct.editCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/collections", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestDeleteCollection_InvalidParam(t *testing.T) {
	mockSvc := new(mockCollectionService)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/collections/delete/abc", nil)
	testutils.SetParam(c, "id", "abc")
	ct.deleteCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/collections", w.Header().Get("Location"))
}

func TestDeleteCollection_DeleteError(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("DeleteByID", uint(5)).Return(errors.New("fail"))
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/collections/delete/5", nil)
	testutils.SetParam(c, "id", "5")
	ct.deleteCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/collections", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestDeleteCollection_Success(t *testing.T) {
	mockSvc := new(mockCollectionService)
	mockSvc.On("DeleteByID", uint(7)).Return(nil)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/collections/delete/7", nil)
	testutils.SetParam(c, "id", "7")
	ct.deleteCollection(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/collections", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}
