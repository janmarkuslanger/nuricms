package field

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
)

func setupTestServer() (*server.Server, *httptest.ResponseRecorder, *testutils.MockFieldService, *testutils.MockCollectionService, *testutils.MockUserService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockField := &testutils.MockFieldService{}
	mockCollection := &testutils.MockCollectionService{}
	mockUser := &testutils.MockUserService{}

	services := &service.Set{
		Field:      mockField,
		Collection: mockCollection,
		User:       mockUser,
	}

	ctrl := NewController(services)
	srv.Handle("GET /fields", ctrl.listFields)
	srv.Handle("GET /fields/create", ctrl.showCreateField)
	srv.Handle("POST /fields/create", ctrl.createField)
	srv.Handle("GET /fields/edit/{id}", ctrl.showEditField)
	srv.Handle("POST /fields/edit/{id}", ctrl.editField)
	srv.Handle("POST /fields/delete/{id}", ctrl.deleteField)

	return srv, rec, mockField, mockCollection, mockUser
}

func Test_listFields(t *testing.T) {
	srv, rec, fieldMock, _, _ := setupTestServer()

	fieldMock.On("List", 1, 10).Return([]model.Field{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/fields", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_showCreateField(t *testing.T) {
	srv, rec, _, collectionMock, _ := setupTestServer()

	collectionMock.On("List", 1, 999999999999999999).Return([]model.Collection{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/fields/create", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_createField(t *testing.T) {
	srv, rec, fieldMock, _, _ := setupTestServer()

	data := dto.FieldData{
		Name:         "title",
		Alias:        "title",
		CollectionID: "1",
		FieldType:    "text",
		IsList:       "false",
		IsRequired:   "true",
		DisplayField: "true",
	}

	fieldMock.On("Create", data).Return(&model.Field{Name: "title"}, nil)

	form := strings.NewReader("name=title&alias=title&collection_id=1&field_type=text&is_list=false&is_required=true&display_field=true")
	req := httptest.NewRequest(http.MethodPost, "/fields/create", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/fields/" {
		t.Errorf("expected redirect to /fields/, got %s", loc)
	}
}

func Test_showEditField(t *testing.T) {
	srv, rec, fieldMock, collectionMock, _ := setupTestServer()

	fieldMock.On("FindByID", uint(1)).Return(&model.Field{Name: "title"}, nil)
	collectionMock.On("List", 1, 999999999999999999).Return([]model.Collection{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/fields/edit/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_editField(t *testing.T) {
	srv, rec, fieldMock, _, _ := setupTestServer()

	data := dto.FieldData{
		Name:         "title",
		Alias:        "title",
		CollectionID: "1",
		FieldType:    "text",
		IsList:       "false",
		IsRequired:   "true",
		DisplayField: "true",
	}

	fieldMock.On("UpdateByID", uint(1), data).Return(&model.Field{Name: "title"}, nil)

	form := strings.NewReader("name=title&alias=title&collection_id=1&field_type=text&is_list=false&is_required=true&display_field=true")
	req := httptest.NewRequest(http.MethodPost, "/fields/edit/1", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/fields/" {
		t.Errorf("expected redirect to /fields/, got %s", loc)
	}
}

func Test_deleteField(t *testing.T) {
	srv, rec, fieldMock, _, _ := setupTestServer()

	fieldMock.On("DeleteByID", uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/fields/delete/1", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/fields/" {
		t.Errorf("expected redirect to /fields/, got %s", loc)
	}
}
