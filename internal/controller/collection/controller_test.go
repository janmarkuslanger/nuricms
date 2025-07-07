package collection

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
)

func setupTestServer() (*server.Server, *httptest.ResponseRecorder, *testutils.MockCollectionService, *testutils.MockUserService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockCollection := &testutils.MockCollectionService{}
	mockUser := &testutils.MockUserService{}

	services := &service.Set{
		Collection: mockCollection,
		User:       mockUser,
	}

	ctrl := NewController(services)
	srv.Handle("GET /collections", ctrl.showCollections)
	srv.Handle("GET /collections/create", ctrl.showCreateCollection)
	srv.Handle("POST /collections/create", ctrl.createCollection)
	srv.Handle("GET /collections/edit/{id}", ctrl.showEditCollection)
	srv.Handle("POST /collections/edit/{id}", ctrl.editCollection)
	srv.Handle("POST /collections/delete/{id}", ctrl.deleteCollection)

	return srv, rec, mockCollection, mockUser
}

func Test_RegisterRoutes(t *testing.T) {
	services := &service.Set{}
	ctrl := NewController(services)

	srv := server.NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/collections", nil)
	ctrl.RegisterRoutes(srv)
	srv.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Errorf("expected different code then 404, got %d", rec.Code)
	}

}

func Test_showCollections(t *testing.T) {
	srv, rec, collectionMock, _ := setupTestServer()
	collectionMock.On("List", 1, 10).Return([]model.Collection{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/collections", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_showCreateCollection(t *testing.T) {
	srv, rec, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/collections/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_createCollection(t *testing.T) {
	srv, rec, collectionMock, _ := setupTestServer()

	data := dto.CollectionData{
		Name:        "Test",
		Alias:       "test",
		Description: "desc",
	}

	collectionMock.
		On("Create", data).
		Return(&model.Collection{Name: "Test"}, nil)

	form := strings.NewReader("name=Test&alias=test&description=desc")
	req := httptest.NewRequest(http.MethodPost, "/collections/create", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/collections" {
		t.Errorf("expected redirect to /collections, got %s", loc)
	}
}

func Test_showEditCollection(t *testing.T) {
	srv, rec, collectionMock, _ := setupTestServer()

	collectionMock.
		On("FindByID", uint(1)).
		Return(&model.Collection{Name: "Test"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/collections/edit/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_editCollection(t *testing.T) {
	srv, rec, collectionMock, _ := setupTestServer()

	data := dto.CollectionData{
		Name:        "Test",
		Alias:       "test",
		Description: "desc",
	}

	collectionMock.
		On("UpdateByID", uint(1), data).
		Return(&model.Collection{Name: "Test"}, nil)

	form := strings.NewReader("name=Test&alias=test&description=desc")
	req := httptest.NewRequest(http.MethodPost, "/collections/edit/1", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/collections/" {
		t.Errorf("expected redirect to /collections/, got %s", loc)
	}
}

func Test_deleteCollection(t *testing.T) {
	srv, rec, collectionMock, _ := setupTestServer()

	collectionMock.
		On("DeleteByID", uint(1)).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/collections/delete/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/collections/" {
		t.Errorf("expected redirect to /collections/, got %s", loc)
	}
}
