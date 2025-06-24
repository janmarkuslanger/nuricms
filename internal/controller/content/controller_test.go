package content

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockController(s *service.Set) *Controller {
	gin.SetMode(gin.TestMode)
	return NewController(s)
}

func TestShowCollections_RendersCorrectly(t *testing.T) {
	mockCollection := &testutils.MockCollectionService{}
	mockCollection.On("List", 1, 10).Return([]model.Collection{
		{Name: "Test"},
	}, int64(1), nil)

	called := false
	originalRender := utils.RenderWithLayout
	utils.RenderWithLayout = func(c *gin.Context, template string, data gin.H, status int) {
		called = true
		assert.Equal(t, "content/collections.tmpl", template)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, 1, data["CurrentPage"])
		assert.Equal(t, 10, data["PageSize"])
	}
	defer func() {
		utils.RenderWithLayout = originalRender
	}()

	s := &service.Set{Collection: mockCollection}
	ct := createMockController(s)

	c, _ := testutils.MakeGETContext("/content/collections?page=1&pageSize=10")

	ct.showCollections(c)

	assert.True(t, called, "RenderWithLayout should have been called")
	mockCollection.AssertExpectations(t)
}

func TestShowCreateContent_RendersCorrectly(t *testing.T) {
	mockField := &testutils.MockFieldService{}
	mockField.On("FindByCollectionID", uint(1)).Return([]model.Field{}, nil)

	mockCol := &testutils.MockCollectionService{}
	mockCol.On("FindByID", uint(1)).Return(&model.Collection{}, nil)

	mockCont := &testutils.MockContentService{}
	mockCont.On("FindContentsWithDisplayContentValue").Return([]model.Content{}, nil)

	mockAsset := &testutils.MockAssetService{}
	mockAsset.On("List", 1, 100000).Return([]model.Asset{}, int64(0), nil)

	s := &service.Set{
		Field:      mockField,
		Collection: mockCol,
		Content:    mockCont,
		Asset:      mockAsset,
	}

	ct := createMockController(s)

	original := utils.RenderWithLayout
	defer func() { utils.RenderWithLayout = original }()

	renderCalled := false
	utils.RenderWithLayout = func(c *gin.Context, template string, data gin.H, code int) {
		renderCalled = true
		assert.Equal(t, "content/create_or_edit.tmpl", template)
		assert.Equal(t, http.StatusOK, code)
	}

	c, _ := testutils.MakeGETContext("/content/collections/1/create")
	testutils.SetParam(c, "id", "1")

	ct.showCreateContent(c)

	assert.True(t, renderCalled, "RenderWithLayout should have been called")

	mockField.AssertExpectations(t)
	mockCol.AssertExpectations(t)
	mockCont.AssertExpectations(t)
	mockAsset.AssertExpectations(t)
}

func TestShowCreateContent_CollectionServiceError(t *testing.T) {
	mockField := &testutils.MockFieldService{}
	mockCollection := &testutils.MockCollectionService{}

	s := &service.Set{
		Field:      mockField,
		Collection: mockCollection,
	}

	ct := createMockController(s)

	c, w := testutils.MakeGETContext("/content/collections/1/create?id=1")
	ct.showCreateContent(c)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
	mockField.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}

func TestCreateContent_Success(t *testing.T) {
	mockContent := new(testutils.MockContentService)
	mockWebhook := new(testutils.MockWebhookService)

	mockContent.On("CreateWithValues", mock.Anything).Return(&model.Content{}, nil)
	mockWebhook.On("Dispatch", string(model.EventContentCreated), nil).Return()

	ct := createMockController(&service.Set{
		Content: mockContent,
		Webhook: mockWebhook,
	})

	c, w := testutils.MakePOSTContext("/content/collections/1/create", gin.H{"title": "Test"})
	testutils.SetParam(c, "id", "1")

	ct.createContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
	mockContent.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestCreateContent_ParseFormError(t *testing.T) {
	mockContent := new(testutils.MockContentService)
	mockWebhook := new(testutils.MockWebhookService)
	ct := createMockController(&service.Set{
		Content: mockContent,
		Webhook: mockWebhook,
	})

	mockContent.On("CreateWithValues", mock.Anything).Return(&model.Content{}, nil)
	mockWebhook.On("Dispatch", string(model.EventContentCreated), nil).Return()

	c, w := testutils.MakeBrokenPOSTContext("/content/collections/1/create")
	testutils.SetParam(c, "id", "1")

	ct.createContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
}

func TestDeleteContent_Valid(t *testing.T) {
	mockContent := new(testutils.MockContentService)
	mockWebhook := new(testutils.MockWebhookService)

	mockContent.On("DeleteByID", uint(42)).Return(nil)
	mockWebhook.On("Dispatch", string(model.EventContentDeleted), nil).Return()

	ct := createMockController(&service.Set{
		Content: mockContent,
		Webhook: mockWebhook,
	})

	c, w := testutils.MakePOSTContext("/content/collections/1/delete/42", nil)
	testutils.SetParam(c, "contentID", "42")

	ct.deleteContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
	mockContent.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestDeleteContent_InvalidParam(t *testing.T) {
	ct := createMockController(&service.Set{})

	c, w := testutils.MakePOSTContext("/content/collections/1/delete/abc", nil)

	ct.deleteContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
}

func TestDeleteContent_ErrorFromService(t *testing.T) {
	mockContent := new(testutils.MockContentService)
	mockWebhook := new(testutils.MockWebhookService)

	mockContent.On("DeleteByID", uint(123)).Return(errors.New("fail"))

	ct := createMockController(&service.Set{
		Content: mockContent,
		Webhook: mockWebhook,
	})

	c, w := testutils.MakePOSTContext("/content/collections/1/delete/123", nil)
	testutils.SetParam(c, "contentID", "123")

	ct.deleteContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
	mockContent.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestEditContent_Success(t *testing.T) {
	mockContent := new(testutils.MockContentService)
	mockWebhook := new(testutils.MockWebhookService)
	ct := createMockController(&service.Set{
		Content: mockContent,
		Webhook: mockWebhook,
	})

	form := gin.H{"title": "Updated Title"}
	c, w := testutils.MakePOSTContext("/content/collections/1/edit/1", form)
	testutils.SetParam(c, "id", "1")
	testutils.SetParam(c, "contentID", "1")

	mockContent.On("EditWithValues", mock.Anything).Return(&model.Content{}, nil)
	mockWebhook.On("Dispatch", string(model.EventContentUpdated), nil).Return()

	ct.editContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
	mockContent.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestEditContent_BadForm(t *testing.T) {
	mockContent := new(testutils.MockContentService)
	mockWebhook := new(testutils.MockWebhookService)
	ct := createMockController(&service.Set{
		Content: mockContent,
		Webhook: mockWebhook,
	})

	c, w := testutils.MakeBrokenPOSTContext("/content/collections/1/edit/1")
	testutils.SetParam(c, "id", "1")
	testutils.SetParam(c, "contentID", "1")

	ct.editContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/content/collections", w.Header().Get("Location"))
}
