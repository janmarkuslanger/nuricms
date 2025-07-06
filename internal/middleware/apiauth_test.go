package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type apimockValidator struct {
	ValidateFunc func(string) error
}

func (m *apimockValidator) Validate(token string) error {
	return m.ValidateFunc(token)
}

func TestApikeyAuth_MissingHeader(t *testing.T) {
	validator := &apimockValidator{}
	handler := ApikeyAuth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "missing X-API-Key")
}

func TestApikeyAuth_InvalidToken(t *testing.T) {
	validator := &apimockValidator{
		ValidateFunc: func(token string) error {
			return errors.New("invalid token")
		},
	}

	handler := ApikeyAuth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "invalid")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid token")
}

func TestApikeyAuth_ValidToken(t *testing.T) {
	validator := &apimockValidator{
		ValidateFunc: func(token string) error {
			return nil
		},
	}

	handler := ApikeyAuth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "valid-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "success", rr.Body.String())
}
