package apikey

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

func setupTestServer() (*server.Server, *httptest.ResponseRecorder, *testutils.MockApikeyService, *testutils.MockUserService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockApikey := &testutils.MockApikeyService{}
	mockUser := &testutils.MockUserService{}

	services := &service.Set{
		Apikey: mockApikey,
		User:   mockUser,
	}

	ctrl := NewController(services)
	srv.Handle("GET /apikeys", ctrl.showApikeys)
	srv.Handle("GET /apikeys/create", ctrl.showCreateApikey)
	srv.Handle("POST /apikeys/create", ctrl.createApikey)
	srv.Handle("GET /apikeys/edit/{id}", ctrl.showEditApikey)
	srv.Handle("POST /apikeys/delete/{id}", ctrl.deleteApikey)

	return srv, rec, mockApikey, mockUser
}

func Test_RegisterRoutes(t *testing.T) {
	services := &service.Set{}
	ctrl := NewController(services)

	srv := server.NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/apikeys", nil)
	ctrl.RegisterRoutes(srv)
	srv.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Errorf("expected different code then 404, got %d", rec.Code)
	}

}

func Test_showApikeys(t *testing.T) {
	srv, rec, apikeyMock, _ := setupTestServer()

	apikeyMock.On("List", 1, 10).Return([]model.Apikey{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/apikeys", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_showCreateApikey(t *testing.T) {
	srv, rec, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/apikeys/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_createApikey(t *testing.T) {
	srv, rec, apikeyMock, _ := setupTestServer()

	apikeyMock.
		On("Create", dto.ApikeyData{Name: "testkey"}).
		Return(&model.Apikey{Name: "testkey"}, nil)

	form := strings.NewReader("name=testkey")
	req := httptest.NewRequest(http.MethodPost, "/apikeys/create", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/apikeys" {
		t.Errorf("expected redirect to /apikeys, got %s", loc)
	}
}

func Test_showEditApikey(t *testing.T) {
	srv, rec, apikeyMock, _ := setupTestServer()

	apikeyMock.
		On("FindByID", uint(1)).
		Return(&model.Apikey{Name: "key"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/apikeys/edit/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_deleteApikey(t *testing.T) {
	srv, rec, apikeyMock, _ := setupTestServer()

	apikeyMock.
		On("DeleteByID", uint(1)).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/apikeys/delete/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/apikeys" {
		t.Errorf("expected redirect to /apikeys, got %s", loc)
	}
}
