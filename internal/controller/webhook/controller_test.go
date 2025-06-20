package webhook

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

func createMockController(webhookSvc service.WebhookService) *Controller {
	gin.SetMode(gin.TestMode)
	svcSet := &service.Set{Webhook: webhookSvc}
	return NewController(svcSet)
}

func TestWebhookController_createWebhook_Success(t *testing.T) {
	svc := new(testutils.MockWebhookService)
	svc.On("Create", "My Webhook", "https://example.com", model.RequestTypePost, mock.MatchedBy(func(events map[model.EventType]bool) bool {
		return events[model.EventContentCreated] == true
	})).Return(&model.Webhook{}, nil)

	ct := createMockController(svc)

	c, w := testutils.MakePOSTContext("/webhooks/create", gin.H{
		"name":           "My Webhook",
		"url":            "https://example.com",
		"request_type":   "POST",
		"ContentCreated": "on",
	})

	ct.createWebhook(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/webhooks", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestWebhookController_createWebhook_Failed(t *testing.T) {
	svc := new(testutils.MockWebhookService)
	svc.On("Create", "Fail Hook", "https://fail", model.RequestTypeGet, mock.Anything).
		Return((*model.Webhook)(nil), errors.New("create failed"))

	ct := createMockController(svc)

	c, w := testutils.MakePOSTContext("/webhooks/create", gin.H{
		"name":         "Fail Hook",
		"url":          "https://fail",
		"request_type": "GET",
	})

	ct.createWebhook(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/webhooks", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestWebhookController_showWebhooks_Success(t *testing.T) {
	svc := new(testutils.MockWebhookService)
	svc.On("List", 1, 10).Return([]model.Webhook{}, int64(0), nil)

	ct := createMockController(svc)

	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakeGETContext("/webhooks?page=1&pageSize=10")
	ct.showWebhooks(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:webhook/index.tmpl")
	svc.AssertExpectations(t)
}

func TestWebhookController_showWebhooks_Failed(t *testing.T) {
	svc := new(testutils.MockWebhookService)
	svc.On("List", 1, 10).Return([]model.Webhook{}, int64(0), errors.New("db error"))

	ct := createMockController(svc)

	c, w := testutils.MakeGETContext("/webhooks?page=1&pageSize=10")
	ct.showWebhooks(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to retrieve webhooks."`)
	svc.AssertExpectations(t)
}
