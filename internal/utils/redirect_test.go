package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestGetParamOrRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		paramValue       string
		expectedID       uint
		expectedOK       bool
		expectRedirect   bool
		expectedLocation string
	}{
		{
			name:           "valid uint param",
			paramValue:     "42",
			expectedID:     42,
			expectedOK:     true,
			expectRedirect: false,
		},
		{
			name:             "invalid uint param",
			paramValue:       "abc",
			expectedID:       0,
			expectedOK:       false,
			expectRedirect:   true,
			expectedLocation: "/fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodGet, "/test/"+tt.paramValue, nil)
			c.Params = gin.Params{
				{Key: "id", Value: tt.paramValue},
			}
			c.Request = req

			id, ok := GetParamOrRedirect(c, "/fallback", "id")

			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedOK, ok)

			if tt.expectRedirect {
				assert.Equal(t, http.StatusSeeOther, w.Code)
				assert.Equal(t, "/fallback", w.Header().Get("Location"))
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}
