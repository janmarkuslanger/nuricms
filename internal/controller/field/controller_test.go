package field

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockController(fieldSvc service.FieldService, collSvc service.CollectionService) *Controller {
	gin.SetMode(gin.TestMode)
	svcSet := &service.Set{Field: fieldSvc, Collection: collSvc}
	return NewController(svcSet)
}

func TestListFields_Success(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockField.On("List", 2, 5).Return([]model.Field{
		{Name: "F1", Alias: "f1"},
		{Name: "F2", Alias: "f2"},
	}, int64(2), nil)
	mockColl := new(testutils.MockCollectionService)
	ct := createMockController(mockField, mockColl)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakeGETContext("/fields?page=2&pageSize=5")
	ct.listFields(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:field/index.tmpl")
	mockField.AssertExpectations(t)
}

func TestListFields_Error(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockField.On("List", 1, 10).Return([]model.Field{}, int64(0), errors.New("fail"))
	mockColl := new(testutils.MockCollectionService)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakeGETContext("/fields")
	ct.listFields(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to retrieve fields."`)
	mockField.AssertExpectations(t)
}

func TestShowCreateField_Success(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	mockColl.On("List", 1, mock.Anything).Return([]model.Collection{
		{Name: "C1", Alias: "c1"},
	}, int64(1), nil)
	ct := createMockController(mockField, mockColl)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakeGETContext("/fields/create")
	ct.showCreateField(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:field/create_or_edit.tmpl")
	mockColl.AssertExpectations(t)
}

func TestShowCreateField_Error(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	mockColl.On("List", 1, mock.Anything).Return([]model.Collection{}, int64(0), errors.New("fail"))
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakeGETContext("/fields/create")
	ct.showCreateField(c)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
	mockColl.AssertExpectations(t)
}

func TestShowEditField_InvalidParam(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakeGETContext("/fields/edit/abc")
	testutils.SetParam(c, "id", "abc")
	ct.showEditField(c)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
}

func TestShowEditField_NotFound(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockField.On("FindByID", uint(5)).Return((*model.Field)(nil), errors.New("not found"))
	mockColl := new(testutils.MockCollectionService)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakeGETContext("/fields/edit/5")
	testutils.SetParam(c, "id", "5")
	ct.showEditField(c)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
	mockField.AssertExpectations(t)
}

func TestShowEditField_CollectionError(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	fieldObj := &model.Field{Alias: "f", CollectionID: 1, FieldType: model.FieldTypeText}
	mockField.On("FindByID", uint(3)).Return(fieldObj, nil)
	mockColl := new(testutils.MockCollectionService)
	mockColl.On("List", 1, mock.Anything).Return([]model.Collection{}, int64(0), errors.New("fail"))
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakeGETContext("/fields/edit/3")
	testutils.SetParam(c, "id", "3")
	ct.showEditField(c)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
	mockField.AssertExpectations(t)
	mockColl.AssertExpectations(t)
}

func TestShowEditField_Success(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	fieldObj := &model.Field{Alias: "f", CollectionID: 1, FieldType: model.FieldTypeText}
	mockField.On("FindByID", uint(4)).Return(fieldObj, nil)
	mockColl := new(testutils.MockCollectionService)
	mockColl.On("List", 1, mock.Anything).Return([]model.Collection{{Name: "C", Alias: "c"}}, int64(1), nil)
	ct := createMockController(mockField, mockColl)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakeGETContext("/fields/edit/4")
	testutils.SetParam(c, "id", "4")
	ct.showEditField(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:field/create_or_edit.tmpl")
	mockField.AssertExpectations(t)
	mockColl.AssertExpectations(t)
}

func TestCreateField_Success(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	data := dto.FieldData{
		Name:         "N",
		Alias:        "a",
		CollectionID: "1",
		FieldType:    "text",
		IsList:       "on",
		IsRequired:   "off",
		DisplayField: "on",
	}
	fieldObj := &model.Field{Name: "N", Alias: "a", CollectionID: 1, FieldType: model.FieldTypeText, IsList: true, DisplayField: true}
	mockField.On("Create", data).Return(fieldObj, nil)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakePOSTContext("/fields/create", gin.H{
		"name":          "N",
		"alias":         "a",
		"collection_id": "1",
		"field_type":    "text",
		"is_list":       "on",
		"is_required":   "off",
		"display_field": "on",
	})
	ct.createField(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
	mockField.AssertExpectations(t)
}

func TestCreateField_Error(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	data := dto.FieldData{
		Name:         "N",
		Alias:        "a",
		CollectionID: "1",
		FieldType:    "text",
		IsList:       "on",
		IsRequired:   "off",
		DisplayField: "on",
	}
	mockField.On("Create", data).Return((*model.Field)(nil), errors.New("fail"))
	ct := createMockController(mockField, mockColl)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakePOSTContext("/fields/create", gin.H{
		"name":          "N",
		"alias":         "a",
		"collection_id": "1",
		"field_type":    "text",
		"is_list":       "on",
		"is_required":   "off",
		"display_field": "on",
	})
	ct.createField(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:field/create_or_edit.tmpl")
	mockField.AssertExpectations(t)
}

func TestEditField_InvalidParam(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakePOSTContext("/fields/edit/abc", gin.H{
		"name":          "N",
		"alias":         "a",
		"collection_id": "1",
		"field_type":    "text",
		"is_list":       "on",
		"is_required":   "off",
		"display_field": "on",
	})
	testutils.SetParam(c, "id", "abc")
	ct.editField(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
}

func TestEditField_UpdateError(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	data := dto.FieldData{
		Name:         "N",
		Alias:        "a",
		CollectionID: "1",
		FieldType:    "text",
		IsList:       "on",
		IsRequired:   "off",
		DisplayField: "on",
	}
	mockField.On("UpdateByID", uint(5), data).Return((*model.Field)(nil), errors.New("fail"))
	ct := createMockController(mockField, mockColl)
	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakePOSTContext("/fields/edit/5", gin.H{
		"name":          "N",
		"alias":         "a",
		"collection_id": "1",
		"field_type":    "text",
		"is_list":       "on",
		"is_required":   "off",
		"display_field": "on",
	})
	testutils.SetParam(c, "id", "5")
	ct.editField(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:field/create_or_edit.tmpl")
	mockField.AssertExpectations(t)
}

func TestEditField_Success(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	data := dto.FieldData{
		Name:         "N",
		Alias:        "a",
		CollectionID: "1",
		FieldType:    "text",
		IsList:       "on",
		IsRequired:   "off",
		DisplayField: "on",
	}
	fieldObj := &model.Field{Name: "N", Alias: "a", CollectionID: 1, FieldType: model.FieldTypeText, IsList: true, DisplayField: true}
	mockField.On("UpdateByID", uint(6), data).Return(fieldObj, nil)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakePOSTContext("/fields/edit/6", gin.H{
		"name":          "N",
		"alias":         "a",
		"collection_id": "1",
		"field_type":    "text",
		"is_list":       "on",
		"is_required":   "off",
		"display_field": "on",
	})
	testutils.SetParam(c, "id", "6")
	ct.editField(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
	mockField.AssertExpectations(t)
}

func TestDeleteField_InvalidParam(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakePOSTContext("/fields/delete/abc", nil)
	testutils.SetParam(c, "id", "abc")
	ct.deleteField(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
}

func TestDeleteField_Success(t *testing.T) {
	mockField := new(testutils.MockFieldService)
	mockColl := new(testutils.MockCollectionService)
	mockField.On("DeleteByID", uint(7)).Return(nil)
	ct := createMockController(mockField, mockColl)

	c, w := testutils.MakePOSTContext("/fields/delete/7", nil)
	testutils.SetParam(c, "id", "7")
	ct.deleteField(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/fields", w.Header().Get("Location"))
	mockField.AssertExpectations(t)
}
