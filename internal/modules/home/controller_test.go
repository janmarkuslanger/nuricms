package home

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils/mockservices"
)

func setupTestServer() (*server.Server, *httptest.ResponseRecorder, *mockservices.MockUserService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockUser := &mockservices.MockUserService{}
	services := &service.Set{
		User: mockUser,
	}

	ctrl := NewController(services)
	srv.Handle("GET /", ctrl.home)

	return srv, rec, mockUser
}

func Test_home(t *testing.T) {
	srv, rec, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
