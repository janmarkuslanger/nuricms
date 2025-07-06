package utils

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultQuery(t *testing.T) {
	t.Run("returns value from query if present", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "page=5",
			},
		}
		result := DefaultQuery(req, "page", "1")
		assert.Equal(t, "5", result)
	})

	t.Run("returns default if query param is missing", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "size=20",
			},
		}
		result := DefaultQuery(req, "page", "1")
		assert.Equal(t, "1", result)
	})

	t.Run("returns default if query param is empty", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "page=",
			},
		}
		result := DefaultQuery(req, "page", "10")
		assert.Equal(t, "10", result)
	})
}
