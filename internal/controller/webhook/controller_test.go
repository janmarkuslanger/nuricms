package webhook

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/mock"
)

func setupWebhookTestServer() (*server.Server, *httptest.ResponseRecorder, *testutils.MockWebhookService) {
	srv := server.NewServer()
	rec := httptest.NewRecorder()

	mockService := &testutils.MockWebhookService{}
	s := &service.Set{
		Webhook: mockService,
		User:    &testutils.MockUserService{},
	}

	ctrl := NewController(s)

	srv.Handle("GET /webhooks", ctrl.showWebhooks)
	srv.Handle("GET /webhooks/create", ctrl.showCreateWebhook)
	srv.Handle("POST /webhooks/create", ctrl.createWebhook)
	srv.Handle("GET /webhooks/edit/{id}", ctrl.showEditWebhook)
	srv.Handle("POST /webhooks/edit/{id}", ctrl.editWebhook)
	srv.Handle("POST /webhooks/delete/{id}", ctrl.deleteWebhook)

	return srv, rec, mockService
}

func Test_RegisterRoutes(t *testing.T) {
	services := &service.Set{}
	ctrl := NewController(services)

	srv := server.NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/webhooks", nil)
	ctrl.RegisterRoutes(srv)
	srv.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Errorf("expected different code then 404, got %d", rec.Code)
	}
}

func Test_showWebhooks(t *testing.T) {
	srv, rec, mockService := setupWebhookTestServer()
	mockService.On("List", mock.Anything, mock.Anything).Return([]model.Webhook{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/webhooks", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_createWebhook(t *testing.T) {
	srv, rec, mockService := setupWebhookTestServer()

	form := url.Values{}
	form.Add("name", "Test Hook")
	form.Add("url", "https://example.com/hook")
	form.Add("request_type", "POST")
	form.Add("content.created", "on")

	events := make(map[string]bool)
	for _, evt := range model.GetWebhookEvents() {
		events[string(evt)] = evt == "content.created"
	}

	mockService.On("Create", dto.WebhookData{
		Name:        "Test Hook",
		Url:         "https://example.com/hook",
		RequestType: "POST",
		Events:      events,
	}).Return(&model.Webhook{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}

func Test_showEditWebhook(t *testing.T) {
	srv, rec, mockService := setupWebhookTestServer()

	mockService.On("FindByID", uint(123)).Return(&model.Webhook{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/webhooks/edit/123", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func Test_editWebhook(t *testing.T) {
	srv, rec, mockService := setupWebhookTestServer()

	form := url.Values{}
	form.Add("name", "Updated Hook")
	form.Add("url", "https://example.com/updated")
	form.Add("request_type", "GET")
	form.Add("content.updated", "on")

	events := make(map[string]bool)
	for _, evt := range model.GetWebhookEvents() {
		events[string(evt)] = form.Get(string(evt)) == "on"
	}

	mockService.On("UpdateByID", uint(123), dto.WebhookData{
		Name:        "Updated Hook",
		Url:         "https://example.com/updated",
		RequestType: "GET",
		Events:      events,
	}).Return(&model.Webhook{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/edit/123", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}

func Test_deleteWebhook(t *testing.T) {
	srv, rec, mockService := setupWebhookTestServer()

	mockService.On("DeleteByID", uint(123)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/delete/123", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}
