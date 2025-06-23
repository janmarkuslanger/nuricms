package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		expectedPage int
		expectedSize int
	}{
		{
			name:         "valid parameters",
			query:        "/?page=3&pageSize=20",
			expectedPage: 3,
			expectedSize: 20,
		},
		{
			name:         "missing parameters",
			query:        "/",
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "invalid parameters",
			query:        "/?page=abc&pageSize=xyz",
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "partial invalid parameters",
			query:        "/?page=5&pageSize=bad",
			expectedPage: 5,
			expectedSize: 10,
		},
		{
			name:         "partial missing parameters",
			query:        "/?pageSize=15",
			expectedPage: 1,
			expectedSize: 15,
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			page, pageSize := ParsePagination(ctx)

			assert.Equal(t, tt.expectedPage, page)
			assert.Equal(t, tt.expectedSize, pageSize)
		})
	}
}
