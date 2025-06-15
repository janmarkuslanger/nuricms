package apikey

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func createMockController(svc service.ApikeyService) *Controller {
	gin.SetMode(gin.TestMode)
	svcSet := &service.Set{Apikey: svc}
	return NewController(svcSet)
}

func TestApiController_ShowApiKeys_Success(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	svc.On("List", 1, 10).Return([]model.Apikey{}, int64(0), nil)
	ct := createMockController(svc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/apikeys?page=1&pageSize=10")
	ct.showApikeys(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:apikey/index.tmpl")
	svc.AssertExpectations(t)
}

func TestApiController_ShowApiKeys_Failed(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	svc.On("List", 1, 10).Return([]model.Apikey{}, int64(0), errors.New("fail"))
	ct := createMockController(svc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/apikeys?page=1&pageSize=10")
	ct.showApikeys(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"fail"`)
	svc.AssertExpectations(t)
}
