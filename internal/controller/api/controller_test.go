package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

func setupApiTestController() (*Controller, *testutils.MockApiService) {
	mockApi := new(testutils.MockApiService)
	services := &service.Set{Api: mockApi}
	return NewController(services), mockApi
}

func TestFindContentById(t *testing.T) {
	r := gin.Default()
	ctrl, mockApi := setupApiTestController()
	r.GET("/api/content/:id", ctrl.findContentById)

	expected := dto.ContentItemResponse{ID: 1}
	mockApi.On("FindContentByID", uint(1)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/api/content/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockApi.AssertExpectations(t)
}

func TestListContents(t *testing.T) {
	r := gin.Default()
	ctrl, mockApi := setupApiTestController()
	r.GET("/api/collections/:alias/content", ctrl.listContents)

	expected := []dto.ContentItemResponse{{ID: 1}}
	mockApi.On("FindContentByCollectionAlias", "blog", 0, 100).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/api/collections/blog/content", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockApi.AssertExpectations(t)
}

func TestListContents_Error(t *testing.T) {
	r := gin.Default()
	ctrl, mockApi := setupApiTestController()
	r.GET("/api/collections/:alias/content", ctrl.listContents)

	err := errors.New("unexpected error")
	mockApi.On("FindContentByCollectionAlias", "blog", 0, 100).Return(nil, err)

	req, _ := http.NewRequest("GET", "/api/collections/blog/content", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockApi.AssertExpectations(t)
}

func TestListContentsByFieldValue(t *testing.T) {
	r := gin.Default()
	ctrl, mockApi := setupApiTestController()
	r.GET("/api/collections/:alias/content/filter", ctrl.listContentsByFieldValue)

	expected := []dto.ContentItemResponse{{ID: 1}}
	mockApi.On("FindContentByCollectionAndFieldValue", "blog", "title", "go", 0, 100).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/api/collections/blog/content/filter?field=title&value=go", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockApi.AssertExpectations(t)
}

func TestListContentsByFieldValue_MissingParams(t *testing.T) {
	r := gin.Default()
	ctrl, _ := setupApiTestController()
	r.GET("/api/collections/:alias/content/filter", ctrl.listContentsByFieldValue)

	req, _ := http.NewRequest("GET", "/api/collections/blog/content/filter", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
