package asset

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/janmarkuslanger/nuricms/testutils/mockservices"
	"github.com/stretchr/testify/mock"
)

func setupAssetTest() (*server.Server, *httptest.ResponseRecorder, *testutils.MockAssetService, *mockservices.MockUserService) {
	s := server.NewServer()
	r := httptest.NewRecorder()
	mockAsset := &testutils.MockAssetService{}
	mockUser := &mockservices.MockUserService{}

	ctrl := NewController(&service.Set{
		Asset: mockAsset,
		User:  mockUser,
	})

	s.Handle("GET /assets", ctrl.showAssets)
	s.Handle("GET /assets/create", ctrl.showCreateAsset)
	s.Handle("POST /assets/create", ctrl.createAsset)
	s.Handle("GET /assets/edit/{id}", ctrl.showEditAsset)
	s.Handle("POST /assets/edit/{id}", ctrl.editAsset)
	s.Handle("POST /assets/delete/{id}", ctrl.deleteAsset)

	return s, r, mockAsset, mockUser
}

func Test_RegisterRoutes(t *testing.T) {
	services := &service.Set{}
	ctrl := NewController(services)

	srv := server.NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/assets", nil)
	ctrl.RegisterRoutes(srv)
	srv.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Errorf("expected different code then 404, got %d", rec.Code)
	}

}

func Test_showAssets(t *testing.T) {
	srv, rec, mockAsset, _ := setupAssetTest()
	mockAsset.On("List", 1, mock.Anything).Return([]model.Asset{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/assets", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func Test_showCreateAsset(t *testing.T) {
	srv, rec, _, _ := setupAssetTest()
	req := httptest.NewRequest(http.MethodGet, "/assets/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_createAsset_success(t *testing.T) {
	srv, rec, mockAsset, _ := setupAssetTest()

	mockAsset.On("UploadFile", mock.Anything, mock.Anything).Return("/uploaded/test.png", nil)
	mockAsset.On("Create", mock.Anything).Return(nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "test-asset")
	part, _ := writer.CreateFormFile("file", "test.png")
	part.Write([]byte("dummy content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/assets/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect, got %d", rec.Code)
	}
}

func Test_deleteAsset_success(t *testing.T) {
	srv, rec, mockAsset, _ := setupAssetTest()
	mockAsset.On("DeleteByID", uint(123)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/assets/delete/123", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect, got %d", rec.Code)
	}
}

func Test_showEditAsset_success(t *testing.T) {
	srv, rec, mockAsset, _ := setupAssetTest()
	mockAsset.On("FindByID", uint(123)).Return(&model.Asset{Name: "Test", Path: "/path"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/assets/edit/123", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func Test_editAsset_withFile(t *testing.T) {
	srv, rec, mockAsset, _ := setupAssetTest()
	mockAsset.On("FindByID", uint(123)).Return(&model.Asset{Name: "Old", Path: "old.png"}, nil)
	mockAsset.On("UploadFile", mock.Anything, mock.Anything).Return("new.png", nil)
	mockAsset.On("Save", mock.Anything).Return(nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Updated Asset")
	part, _ := writer.CreateFormFile("file", "update.png")
	part.Write([]byte("updated content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/assets/edit/123", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect, got %d", rec.Code)
	}
}

func Test_editAsset_withoutFile(t *testing.T) {
	srv, rec, mockAsset, _ := setupAssetTest()
	mockAsset.On("FindByID", uint(123)).Return(&model.Asset{Name: "Old", Path: "old.png"}, nil)
	mockAsset.On("Save", mock.Anything).Return(nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Updated Name")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/assets/edit/123", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect, got %d", rec.Code)
	}
}
