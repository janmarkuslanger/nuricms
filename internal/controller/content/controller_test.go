// Filename: internal/controller/content/controller_test.go

package content

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockController(s *service.Set) *Controller {
	gin.SetMode(gin.TestMode)
	return NewController(s)
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
	mockWebhook.On("Dispatch", string(model.EventContentDeleted), nil).Return()

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
