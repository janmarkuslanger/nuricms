package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/testutils"
)

func newTestWebhookServiceWithClient(t *testing.T, client *http.Client) *webhookService {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	return &webhookService{
		repos:      repos,
		httpClient: client,
	}
}

func newTestWebhookService(t *testing.T) WebhookService {
	return NewWebhookService(repository.NewSet(testutils.SetupTestDB(t)))
}

func TestWebhookService_Create_Success(t *testing.T) {
	svc := newTestWebhookService(t)

	hook, err := svc.Create(dto.WebhookData{
		Name:        "My Hook",
		Url:         "https://example.com",
		RequestType: "POST",
		Events: map[string]bool{
			"create": true,
			"update": false,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, hook)
	assert.Equal(t, "My Hook", hook.Name)
	assert.Equal(t, "https://example.com", hook.Url)
	assert.Contains(t, hook.Events, "create")
	assert.NotContains(t, hook.Events, "update")
}

func TestWebhookService_Create_ValidationError(t *testing.T) {
	svc := newTestWebhookService(t)

	_, err := svc.Create(dto.WebhookData{})
	assert.EqualError(t, err, "no name given")

	_, err = svc.Create(dto.WebhookData{Name: "x"})
	assert.EqualError(t, err, "no url given")

	_, err = svc.Create(dto.WebhookData{Name: "x", Url: "y"})
	assert.EqualError(t, err, "no request type given")
}

func TestWebhookService_UpdateByID(t *testing.T) {
	svc := newTestWebhookService(t)

	hook, _ := svc.Create(dto.WebhookData{
		Name:        "Original",
		Url:         "https://original.com",
		RequestType: "POST",
		Events:      map[string]bool{"create": true},
	})

	updated, err := svc.UpdateByID(hook.ID, dto.WebhookData{
		Name:        "Updated",
		Url:         "https://updated.com",
		RequestType: "POST",
		Events:      map[string]bool{"delete": true},
	})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
	assert.Equal(t, "https://updated.com", updated.Url)
	assert.Contains(t, updated.Events, "delete")
}

func TestWebhookService_UpdateByID_ValidationError(t *testing.T) {
	svc := newTestWebhookService(t)
	_, err := svc.UpdateByID(999, dto.WebhookData{})
	assert.EqualError(t, err, "no name given")
}

func TestWebhookService_List_Find_Save_Delete(t *testing.T) {
	svc := newTestWebhookService(t)

	h, _ := svc.Create(dto.WebhookData{
		Name:        "Hook A",
		Url:         "https://a.com",
		RequestType: "POST",
		Events:      map[string]bool{"create": true},
	})

	list, total, err := svc.List(1, 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int64(len(list)), int64(1))
	assert.GreaterOrEqual(t, total, int64(1))

	found, err := svc.FindByID(h.ID)
	assert.NoError(t, err)
	assert.Equal(t, h.ID, found.ID)

	found.Active = false
	err = svc.Save(found)
	assert.NoError(t, err)

	err = svc.DeleteByID(h.ID)
	assert.NoError(t, err)

	_, err = svc.FindByID(h.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestWebhookService_Dispatch(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		defer r.Body.Close()
		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)
		assert.Equal(t, "value", payload["key"])
	}))
	defer server.Close()

	svc := newTestWebhookServiceWithClient(t, server.Client())
	hook, _ := svc.Create(dto.WebhookData{
		Name:        "DispatchHook",
		Url:         server.URL,
		RequestType: "POST",
		Events:      map[string]bool{"dispatch": true},
	})

	svc.Dispatch("dispatch", map[string]string{"key": "value"})
	time.Sleep(100 * time.Millisecond)

	assert.True(t, called, "Dispatch should have triggered HTTP call")
	assert.NotZero(t, hook.ID)
}
