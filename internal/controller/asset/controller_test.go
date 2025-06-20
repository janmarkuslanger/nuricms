package asset

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockController(svc service.AssetService) *Controller {
	gin.SetMode(gin.TestMode)
	svcSet := &service.Set{Asset: svc}
	return NewController(svcSet)
}

func TestShowAssets_Success(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	assets := []model.Asset{
		{Name: "A", Path: "/path/a"},
		{Name: "B", Path: "/path/b"},
	}
	mockSvc.On("List", 2, 5).Return(assets, int64(2), nil)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/assets?page=2&pageSize=5")
	ct.showAssets(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:asset/index.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestShowAssets_InvalidQueryParams(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	mockSvc.On("List", 1, 10).Return([]model.Asset{}, int64(0), nil)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/assets?page=abc&pageSize=xyz")
	ct.showAssets(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:asset/index.tmpl")
	mockSvc.AssertExpectations(t)
}

func TestShowCreateAsset(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/assets/create")
	ct.showCreateAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:asset/create_or_edit.tmpl")
}

func TestDeleteAsset_InvalidParam(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	mockSvc.On("DeleteByID", uint(1)).Return(errors.New("no"))
	c, w := testutils.MakePOSTContext("/assets/delete/abc", nil)
	testutils.SetParam(c, "id", "abc")
	ct.deleteAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertNotCalled(t, "DeleteByID")
}

func TestDeleteAsset_Error(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	mockSvc.On("DeleteByID", uint(5)).Return(errors.New("Fai"))
	c, w := testutils.MakePOSTContext("/assets/delete/5", nil)
	testutils.SetParam(c, "id", "5")
	ct.deleteAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestDeleteAsset_Success(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	mockSvc.On("DeleteByID", uint(7)).Return(nil)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/assets/delete/7", nil)
	testutils.SetParam(c, "id", "7")
	ct.deleteAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestShowEditAsset_InvalidParam(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	mockSvc.On("FindByID", uint(3)).Return((*model.Asset)(nil), nil)
	c, w := testutils.MakeGETContext("/assets/edit/abc")
	testutils.SetParam(c, "id", "abc")
	ct.showEditAsset(c)
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
}

func TestShowEditAsset_NotFound(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	mockSvc.On("FindByID", uint(3)).Return((*model.Asset)(nil), assert.AnError)
	ct := createMockController(mockSvc)
	c, w := testutils.MakeGETContext("/assets/edit/3")
	testutils.SetParam(c, "id", "3")
	ct.showEditAsset(c)
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestShowEditAsset_Success(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	asset := &model.Asset{Name: "N", Path: "/p"}
	mockSvc.On("FindByID", uint(4)).Return(asset, nil)
	ct := createMockController(mockSvc)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()
	c, w := testutils.MakeGETContext("/assets/edit/4")
	testutils.SetParam(c, "id", "4")
	ct.showEditAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:asset/create_or_edit.tmpl")
	mockSvc.AssertExpectations(t)
}

func makeMultipartRequest(path, fieldName, fileName, fileContent string, formFields map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fieldName, fileName)
	io.Copy(part, bytes.NewBufferString(fileContent))
	for k, v := range formFields {
		writer.WriteField(k, v)
	}
	writer.Close()
	req := httptest.NewRequest("POST", path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func TestCreateAsset_NoFile(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	body := "name=Test"
	req := httptest.NewRequest("POST", "/assets/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	ct.createAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertNotCalled(t, "UploadFile")
	mockSvc.AssertNotCalled(t, "Create")
}

func TestCreateAsset_UploadError(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	c, w := makeMultipartRequest("/assets/create", "file", "f.txt", "data", map[string]string{"name": "TestName"})
	mockSvc.On("UploadFile", mock.Anything, mock.Anything).Return("", assert.AnError)
	ct.createAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestCreateAsset_Success(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	c, w := makeMultipartRequest("/assets/create", "file", "f.txt", "content", map[string]string{"name": "MyName"})
	mockSvc.On("UploadFile", mock.Anything, mock.Anything).Return("/tmp/f.txt", nil)
	mockSvc.On("Create", &model.Asset{Name: "MyName", Path: "/tmp/f.txt"}).Return(nil)
	ct.createAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestEditAsset_InvalidParam(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/assets/edit/abc", nil)
	testutils.SetParam(c, "id", "abc")
	ct.editAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
}

func TestEditAsset_FindError(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	mockSvc.On("FindByID", uint(5)).Return((*model.Asset)(nil), assert.AnError)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/assets/edit/5", nil)
	testutils.SetParam(c, "id", "5")
	ct.editAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestEditAsset_WithFile_Success(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	asset := &model.Asset{Name: "Old", Path: "/old"}
	mockSvc.On("FindByID", uint(6)).Return(asset, nil)
	c, w := makeMultipartRequest("/assets/edit/6", "file", "new.txt", "newdata", map[string]string{"name": "NewName"})
	testutils.SetParam(c, "id", "6")
	mockSvc.On("UploadFile", mock.Anything, mock.Anything).Return("/new/path", nil)
	mockSvc.On("Save", &model.Asset{Name: "NewName", Path: "/new/path"}).Return(nil)
	ct := createMockController(mockSvc)
	ct.editAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}

func TestEditAsset_WithoutFile_Success(t *testing.T) {
	mockSvc := new(testutils.MockAssetService)
	asset := &model.Asset{Name: "OldName", Path: "/path"}
	mockSvc.On("FindByID", uint(8)).Return(asset, nil)
	ct := createMockController(mockSvc)
	c, w := testutils.MakePOSTContext("/assets/edit/8", gin.H{"name": "Updated"})
	testutils.SetParam(c, "id", "8")
	mockSvc.On("Save", &model.Asset{Name: "Updated", Path: "/path"}).Return(nil)
	ct.editAsset(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/assets", w.Header().Get("Location"))
	mockSvc.AssertExpectations(t)
}
