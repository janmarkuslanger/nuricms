package apikey

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
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

func TestApiController_ShowEditApiKeys_Success(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	svc.On("FindByID", uint(1)).Return(&model.Apikey{}, nil)
	ct := createMockController(svc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/apikeys/1")
	testutils.SetParam(c, "id", "1")
	ct.showEditApikey(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:apikey/create_or_edit.tmpl")
	svc.AssertExpectations(t)
}

func TestApiController_CreateApiKey_Success(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	svc.On("CreateToken", "TestKey", time.Duration(0)).Return("hello_world", nil)
	ct := createMockController(svc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakePOSTContext("/apikeys/create", gin.H{
		"name": "TestKey",
	})
	ct.createApikey(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/apikeys", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestApiController_CreateApiKey_EmptyName(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	svc.On("CreateToken", "", time.Duration(0)).Return("", errors.New("Is error"))
	ct := createMockController(svc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakePOSTContext("/apikeys/create", gin.H{
		"name": "",
	})
	ct.createApikey(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/apikeys/create", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestShowEditApikey_InvalidParam(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	ct := createMockController(svc)
	c, w := testutils.MakeGETContext("/apikeys/edit/abc")
	testutils.SetParam(c, "id", "abc")
	ct.showEditApikey(c)
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/apikeys", w.Header().Get("Location"))
}

func TestShowEditApikey_NotFound(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	svc.On("FindByID", uint(5)).Return((*model.Apikey)(nil), errors.New("not found"))
	ct := createMockController(svc)
	c, w := testutils.MakeGETContext("/apikeys/edit/5")
	testutils.SetParam(c, "id", "5")
	ct.showEditApikey(c)
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/apikeys", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestShowEditApikey_Success(t *testing.T) {
	svc := new(testutils.MockApikeyService)
	ap := &model.Apikey{Name: "Key", Token: "tok"}
	svc.On("FindByID", uint(1)).Return(ap, nil)
	ct := createMockController(svc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/apikeys/edit/1")
	testutils.SetParam(c, "id", "1")
	ct.showEditApikey(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:apikey/create_or_edit.tmpl")
	svc.AssertExpectations(t)
}
