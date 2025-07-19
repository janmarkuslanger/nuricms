package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
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
	srv.Handle("GET /login", ctrl.showLogin)
	srv.Handle("POST /login", ctrl.login)
	srv.Handle("GET /user", ctrl.showUser)
	srv.Handle("GET /user/create", ctrl.showCreateUser)
	srv.Handle("POST /user/create", ctrl.createUser)
	srv.Handle("GET /user/edit/{id}", ctrl.showEditUser)
	srv.Handle("POST /user/edit/{id}", ctrl.editUser)
	srv.Handle("POST /user/delete/{id}", ctrl.deleteUser)

	return srv, rec, mockUser
}

func Test_RegisterRoutes(t *testing.T) {
	services := &service.Set{}
	ctrl := NewController(services)

	srv := server.NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	ctrl.RegisterRoutes(srv)
	srv.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Errorf("expected different code then 404, got %d", rec.Code)
	}
}

func Test_showLogin(t *testing.T) {
	srv, rec, _ := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_login(t *testing.T) {
	srv, rec, mockUser := setupTestServer()

	mockUser.
		On("LoginUser", "test@example.com", "securepassword").
		Return("mocktoken", nil)

	form := strings.NewReader("email=test@example.com&password=securepassword")
	req := httptest.NewRequest(http.MethodPost, "/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect (303), got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/" {
		t.Errorf("expected redirect to /, got %s", loc)
	}
}

func Test_showUser(t *testing.T) {
	srv, rec, mockUser := setupTestServer()
	mockUser.On("List", 1, 10).Return([]model.User{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_showCreateUser(t *testing.T) {
	srv, rec, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/user/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_createUser(t *testing.T) {
	srv, rec, mockUser := setupTestServer()

	mockUser.
		On("Create", dto.UserData{
			Email:    "test@example.com",
			Password: "securepassword",
			Role:     "admin",
		}).
		Return(&model.User{Email: "test@example.com"}, nil)

	form := strings.NewReader("email=test@example.com&password=securepassword&role=admin")
	req := httptest.NewRequest(http.MethodPost, "/user/create", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/user" {
		t.Errorf("expected redirect to /user, got %s", loc)
	}
}

func Test_showEditUser(t *testing.T) {
	srv, rec, mockUser := setupTestServer()

	mockUser.
		On("FindByID", uint(1)).
		Return(&model.User{Email: "user@example.com"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/user/edit/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_editUser(t *testing.T) {
	srv, rec, mockUser := setupTestServer()

	mockUser.
		On("UpdateByID", uint(1), dto.UserData{
			Email:    "updated@example.com",
			Password: "newpass",
			Role:     "editor",
		}).
		Return(&model.User{Email: "updated@example.com"}, nil)

	form := strings.NewReader("email=updated@example.com&password=newpass&role=editor")
	req := httptest.NewRequest(http.MethodPost, "/user/edit/1", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/user" {
		t.Errorf("expected redirect to /user, got %s", loc)
	}
}

func Test_deleteUser(t *testing.T) {
	srv, rec, mockUser := setupTestServer()

	mockUser.
		On("DeleteByID", uint(1)).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/user/delete/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/user" {
		t.Errorf("expected redirect to /user, got %s", loc)
	}
}
