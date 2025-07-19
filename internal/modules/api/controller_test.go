package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
)

func setupTestServer() (*server.Server, *httptest.ResponseRecorder, *testutils.MockApiService, *testutils.MockApikeyService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockApi := &testutils.MockApiService{}
	mockApikey := &testutils.MockApikeyService{}

	services := &service.Set{
		Api:    mockApi,
		Apikey: mockApikey,
	}

	ctrl := NewController(services)

	srv.Handle("GET /api/collections/{alias}/content", ctrl.listContents)
	srv.Handle("GET /api/content/{id}", ctrl.findContentById)
	srv.Handle("GET /api/collections/{alias}/content/filter", ctrl.listContentsByFieldValue)

	return srv, rec, mockApi, mockApikey
}

func Test_RegisterRoutes(t *testing.T) {
	services := &service.Set{}
	ctrl := NewController(services)

	srv := server.NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/collections/lol/content", nil)
	ctrl.RegisterRoutes(srv)
	srv.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Errorf("expected different code then 404, got %d", rec.Code)
	}

}

func Test_findContentById(t *testing.T) {
	srv, rec, mockApi, _ := setupTestServer()

	mockApi.
		On("FindContentByID", uint(1)).
		Return(dto.ContentItemResponse{
			ID: 1,
			Collection: dto.CollectionResponse{
				ID:    2,
				Name:  "Example",
				Alias: "example",
			},
			Values: map[string]any{},
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/content/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var raw struct {
		Success    bool                    `json:"success"`
		Data       dto.ContentItemResponse `json:"data"`
		Meta       *dto.MetaData           `json:"meta"`
		Pagination *dto.Pagination         `json:"pagination"`
		Error      *dto.ErrorDetail        `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !raw.Success {
		t.Errorf("expected success to be true, got false")
	}

	if raw.Data.ID != 1 {
		t.Errorf("expected ID 1, got %d", raw.Data.ID)
	}

	if raw.Data.Collection.Alias != "example" {
		t.Errorf("expected collection alias 'example', got '%s'", raw.Data.Collection.Alias)
	}
}

func Test_listContents(t *testing.T) {
	srv, rec, mockApi, _ := setupTestServer()

	mockApi.
		On("FindContentByCollectionAlias", "news", 0, 100).
		Return([]dto.ContentItemResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/collections/news/content", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp dto.ApiResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if !resp.Success || len(resp.Data.([]interface{})) != 1 {
		t.Errorf("unexpected response data: %+v", resp)
	}
}

func Test_listContentsByFieldValue(t *testing.T) {
	srv, rec, mockApi, _ := setupTestServer()

	mockApi.
		On("FindContentByCollectionAndFieldValue", "blog", "title", "Golang", 0, 100).
		Return([]dto.ContentItemResponse{{ID: 42}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/collections/blog/content/filter?field=title&value=Golang", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp dto.ApiResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if !resp.Success || len(resp.Data.([]interface{})) != 1 {
		t.Errorf("unexpected response data: %+v", resp)
	}
}
