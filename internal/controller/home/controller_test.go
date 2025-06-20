package home

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

func createMockController() *Controller {
	gin.SetMode(gin.TestMode)
	return NewController(&service.Set{User: nil})
}

func TestHome_Success(t *testing.T) {
	ct := createMockController()
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakeGETContext("/")
	ct.home(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:home/home.tmpl")
}
