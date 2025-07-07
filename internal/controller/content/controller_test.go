package content

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/mock"
)

func setup(t *testing.T) (*server.Server, *httptest.ResponseRecorder, *testutils.MockCollectionService, *testutils.MockContentService, *testutils.MockFieldService, *testutils.MockAssetService, *testutils.MockWebhookService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockColl := &testutils.MockCollectionService{}
	mockCont := &testutils.MockContentService{}
	mockField := &testutils.MockFieldService{}
	mockAsset := &testutils.MockAssetService{}
	mockWebhook := &testutils.MockWebhookService{}
	mockUser := &testutils.MockUserService{}

	services := &service.Set{
		Collection: mockColl,
		Content:    mockCont,
		Field:      mockField,
		Asset:      mockAsset,
		Webhook:    mockWebhook,
		User:       mockUser,
	}

	ctrl := NewController(services)

	srv.Handle("GET /content/collections", ctrl.showCollections)
	srv.Handle("GET /content/collections/{id}/show", ctrl.listContent)
	srv.Handle("GET /content/collections/{id}/create", ctrl.showCreateContent)
	srv.Handle("POST /content/collections/{id}/create", ctrl.createContent)
	srv.Handle("GET /content/collections/{id}/edit/{contentID}", ctrl.showEditContent)
	srv.Handle("POST /content/collections/{id}/edit/{contentID}", ctrl.editContent)
	srv.Handle("POST /content/collections/{id}/delete/{contentID}", ctrl.deleteContent)

	return srv, rec, mockColl, mockCont, mockField, mockAsset, mockWebhook
}

func Test_RegisterRoutes(t *testing.T) {
	services := &service.Set{}
	ctrl := NewController(services)

	srv := server.NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/content/collections", nil)
	ctrl.RegisterRoutes(srv)
	srv.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Errorf("expected different code then 404, got %d", rec.Code)
	}
}

func Test_showCollections(t *testing.T) {
	srv, rec, mockColl, _, _, _, _ := setup(t)

	mockColl.On("List", 1, mock.Anything).Return([]model.Collection{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_createContent_success(t *testing.T) {
	srv, rec, _, mockCont, _, _, mockWebhook := setup(t)

	form := url.Values{}
	form.Add("field_1", "value")

	mockCont.On("CreateWithValues", mock.Anything).Return(&model.Content{}, nil)
	mockWebhook.On("Dispatch", string(model.EventContentCreated), mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/content/collections/1/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}

func Test_editContent_success(t *testing.T) {
	srv, rec, _, mockCont, _, _, mockWebhook := setup(t)

	form := url.Values{}
	form.Add("field_1", "value")

	mockCont.On("EditWithValues", mock.Anything).Return(&model.Content{}, nil)
	mockWebhook.On("Dispatch", string(model.EventContentUpdated), mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/content/collections/1/edit/2", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}

func Test_deleteContent_success(t *testing.T) {
	srv, rec, _, mockCont, _, _, mockWebhook := setup(t)

	mockCont.On("DeleteByID", uint(2)).Return(nil)
	mockWebhook.On("Dispatch", string(model.EventContentDeleted), mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/content/collections/1/delete/2", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}
