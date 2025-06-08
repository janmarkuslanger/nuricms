package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeApikeyService struct {
	validKeys map[string]error
}

func (f *fakeApikeyService) Validate(token string) error {
	if err, ok := f.validKeys[token]; ok {
		return err
	}
	return errors.New("invalid API key")
}

func TestApikeyAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeSvc := &fakeApikeyService{
		validKeys: map[string]error{
			"valid-key": nil,
			"bad-key":   errors.New("key revoked"),
		},
	}

	tests := []struct {
		name           string
		setHeader      bool
		headerValue    string
		expectedStatus int
		expectedBody   string
		callNext       bool
	}{
		{
			name:           "Missing X-API-Key header → 401",
			setHeader:      false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"missing X-API-Key header"}`,
			callNext:       false,
		},
		{
			name:           "Invalid key → 401",
			setHeader:      true,
			headerValue:    "unknown-key",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid API key"}`,
			callNext:       false,
		},
		{
			name:           "Revoked key → 401",
			setHeader:      true,
			headerValue:    "bad-key",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"key revoked"}`,
			callNext:       false,
		},
		{
			name:           "Valid key → Next aufgerufen",
			setHeader:      true,
			headerValue:    "valid-key",
			expectedStatus: http.StatusOK,
			callNext:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ApikeyAuth(fakeSvc))
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tc.setHeader {
				req.Header.Set("X-API-Key", tc.headerValue)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.callNext {
				assert.Equal(t, "OK", w.Body.String())
			} else {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}
