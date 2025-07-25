package content

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/janmarkuslanger/nuricms/testutils/mockservices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setup(t *testing.T) (*server.Server, *httptest.ResponseRecorder, *testutils.MockCollectionService, *testutils.MockContentService, *testutils.MockFieldService, *testutils.MockAssetService, *testutils.MockWebhookService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockColl := &testutils.MockCollectionService{}
	mockCont := &testutils.MockContentService{}
	mockField := &testutils.MockFieldService{}
	mockAsset := &testutils.MockAssetService{}
	mockWebhook := &testutils.MockWebhookService{}
	mockUser := &mockservices.MockUserService{}

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

func Test_showCreateContent_success(t *testing.T) {
	srv, rec, mockColl, mockContent, mockField, mockAsset, _ := setup(t)

	collectionID := uint(1)

	mockField.On("FindByCollectionID", collectionID).Return([]model.Field{
		{Alias: "title"},
	}, nil)

	mockColl.On("FindByID", collectionID).Return(&model.Collection{}, nil)

	mockContent.On("FindContentsWithDisplayContentValue").Return([]model.Content{}, nil)

	mockAsset.On("List", 1, 100000).Return([]model.Asset{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections/1/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}
}

func Test_showCreateContent_paramredirect(t *testing.T) {
	srv, rec, mockColl, mockContent, mockField, mockAsset, _ := setup(t)

	collectionID := uint(1)

	mockField.On("FindByCollectionID", collectionID).Return([]model.Field{
		{Alias: "title"},
	}, nil)

	mockColl.On("FindByID", collectionID).Return(&model.Collection{}, nil)

	mockContent.On("FindContentsWithDisplayContentValue").Return([]model.Content{}, nil)

	mockAsset.On("List", 1, 100000).Return([]model.Asset{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections/asdfsadf/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 OK, got %d", rec.Code)
	}
}

func Test_showCreateContent_notfound(t *testing.T) {
	srv, rec, mockColl, mockContent, mockField, mockAsset, _ := setup(t)

	collectionID := uint(1)

	mockField.On("FindByCollectionID", collectionID).Return([]model.Field{
		{Alias: "title"},
	}, nil)

	mockColl.On("FindByID", collectionID).Return(&model.Collection{}, errors.New("asd"))

	mockContent.On("FindContentsWithDisplayContentValue").Return([]model.Content{}, nil)

	mockAsset.On("List", 1, 100000).Return([]model.Asset{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections/asdfsadf/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 OK, got %d", rec.Code)
	}
}

func Test_listContent_success(t *testing.T) {
	srv, rec, _, mockContent, mockField, _, _ := setup(t)

	collectionID := uint(1)
	page := 1
	pageSize := 10

	mockContent.On("FindDisplayValueByCollectionID", collectionID, page, pageSize).
		Return([]model.Content{}, int64(0), nil)

	mockField.On("FindDisplayFieldsByCollectionID", collectionID).
		Return([]model.Field{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections/1/show", nil)
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_showEditContent_success(t *testing.T) {
	srv, rec, mockColl, mockContent, _, mockAsset, _ := setup(t)

	collectionID := uint(1)
	contentID := uint(42)

	mockContent.On("FindByID", contentID).Return(&model.Content{Model: gorm.Model{ID: contentID}}, nil)
	mockColl.On("FindByID", collectionID).Return(&model.Collection{Model: gorm.Model{ID: collectionID}}, nil)
	mockContent.On("FindContentsWithDisplayContentValue").Return([]model.Content{}, nil)
	mockAsset.On("List", 1, 100000).Return([]model.Asset{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections/1/edit/42", nil)
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_showEditContent_paramredirect(t *testing.T) {
	srv, rec, mockColl, mockContent, _, mockAsset, _ := setup(t)

	collectionID := uint(1)
	contentID := uint(42)

	mockContent.On("FindByID", contentID).Return(&model.Content{Model: gorm.Model{ID: contentID}}, nil)
	mockColl.On("FindByID", collectionID).Return(&model.Collection{Model: gorm.Model{ID: collectionID}}, nil)
	mockContent.On("FindContentsWithDisplayContentValue").Return([]model.Content{}, nil)
	mockAsset.On("List", 1, 100000).Return([]model.Asset{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections/1/edit/qwe", nil)
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
}

func Test_showEditContent_notfound(t *testing.T) {
	srv, rec, mockColl, mockContent, _, mockAsset, _ := setup(t)

	collectionID := uint(1)
	contentID := uint(42)

	mockContent.On("FindByID", contentID).Return(&model.Content{Model: gorm.Model{ID: contentID}}, errors.New("no no"))
	mockColl.On("FindByID", collectionID).Return(&model.Collection{Model: gorm.Model{ID: collectionID}}, nil)
	mockContent.On("FindContentsWithDisplayContentValue").Return([]model.Content{}, nil)
	mockAsset.On("List", 1, 100000).Return([]model.Asset{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/content/collections/1/edit/qwe", nil)
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
}
