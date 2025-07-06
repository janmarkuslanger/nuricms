package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockValidator struct {
	mock.Mock
}

func (m *mockValidator) ValidateJWT(token string) (uint, string, model.Role, error) {
	args := m.Called(token)
	return args.Get(0).(uint), args.String(1), args.Get(2).(model.Role), args.Error(3)
}

func TestUserauth_WithCookie_Success(t *testing.T) {
	validator := new(mockValidator)
	validator.On("ValidateJWT", "validtoken").Return(uint(1), "user@example.com", model.RoleAdmin, nil)

	handler := Userauth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, uint(1), r.Context().Value(UserIDKey))
		assert.Equal(t, "user@example.com", r.Context().Value(UserEmailKey))
		assert.Equal(t, model.RoleAdmin, r.Context().Value(UserRoleKey))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "validtoken"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserauth_WithHeader_Success(t *testing.T) {
	validator := new(mockValidator)
	validator.On("ValidateJWT", "validtoken").Return(uint(2), "header@example.com", model.RoleEditor, nil)

	handler := Userauth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, uint(2), r.Context().Value(UserIDKey))
		assert.Equal(t, "header@example.com", r.Context().Value(UserEmailKey))
		assert.Equal(t, model.RoleEditor, r.Context().Value(UserRoleKey))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserauth_MissingToken(t *testing.T) {
	validator := new(mockValidator)

	handler := Userauth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not reach handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestRoleauth_AccessGranted(t *testing.T) {
	handler := Roleauth(model.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, UserRoleKey, model.RoleAdmin)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleauth_AccessDenied(t *testing.T) {
	handler := Roleauth(model.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be allowed")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, UserRoleKey, model.RoleEditor)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "insufficient permissions"))
}

func TestRoleauth_NoRoleSet(t *testing.T) {
	handler := Roleauth(model.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be allowed")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "missing role"))
}
