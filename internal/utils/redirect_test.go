package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetParamOrRedirect_Valid(t *testing.T) {
	req := httptest.NewRequest("GET", "/items/99", nil)
	req.SetPathValue("id", "99")

	rec := httptest.NewRecorder()
	ctx := server.Context{
		Request: req,
		Writer:  rec,
	}

	id, ok := utils.GetParamOrRedirect(ctx, "/redirect", "id")

	assert.Equal(t, ok, true)
	assert.Equal(t, uint(99), id)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetParamOrRedirect_InvalidParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/items/abc", nil)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()
	ctx := server.Context{
		Request: req,
		Writer:  rec,
	}

	id, ok := utils.GetParamOrRedirect(ctx, "/fallback", "id")

	assert.Equal(t, ok, false)
	assert.Equal(t, uint(0), id)
	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/fallback", rec.Header().Get("Location"))
}
