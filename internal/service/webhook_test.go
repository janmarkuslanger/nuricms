package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/testutils"
)

func TestWebhookService_Create_List_Find_Save_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svcIface := NewWebhookService(repos)
	ws := svcIface.(*webhookService)

	events := map[model.EventType]bool{
		model.EventContentCreated: true,
		model.EventContentDeleted: false,
		model.EventContentUpdated: true,
	}
	wh, err := ws.Create("name1", "http://example.com", model.RequestTypePost, events)
	assert.NoError(t, err)
	assert.NotZero(t, wh.ID)
	assert.Contains(t, wh.Events, string(model.EventContentCreated)+",")
	assert.Contains(t, wh.Events, string(model.EventContentUpdated)+",")
	assert.NotContains(t, wh.Events, string(model.EventContentDeleted)+",")

	wh2, err := ws.Create("name2", "http://example.org", model.RequestTypeGet, map[model.EventType]bool{"e": true})
	assert.NoError(t, err)

	list, total, err := ws.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	ids := []uint{list[0].ID, list[1].ID}
	assert.Contains(t, ids, wh.ID)
	assert.Contains(t, ids, wh2.ID)

	found, err := ws.FindByID(wh.ID)
	assert.NoError(t, err)
	assert.Equal(t, wh.ID, found.ID)

	found.Url = "http://changed"
	err = ws.Save(found)
	assert.NoError(t, err)
	reloaded, err := ws.FindByID(wh.ID)
	assert.NoError(t, err)
	assert.Equal(t, "http://changed", reloaded.Url)

	err = ws.DeleteByID(wh2.ID)
	assert.NoError(t, err)
	_, err = repos.Webhook.FindByID(wh2.ID)
	assert.Error(t, err)
}

func TestWebhookService_Dispatch_NoHooks(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svcIface := NewWebhookService(repos)
	ws := svcIface.(*webhookService)
	ws.Dispatch("noevent", map[string]string{"k": "v"})
	time.Sleep(10 * time.Millisecond)
}

func TestWebhookService_Dispatch_SendsRequest(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svcIface := NewWebhookService(repos)
	ws := svcIface.(*webhookService)

	var wg sync.WaitGroup
	payloadCh := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		body, _ := io.ReadAll(r.Body)
		payloadCh <- body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ws.httpClient = server.Client()

	evs := map[model.EventType]bool{"testev": true}
	wh, err := ws.Create("w1", server.URL, model.RequestTypePost, evs)
	assert.NoError(t, err)
	assert.NotZero(t, wh.ID)

	payload := map[string]interface{}{"foo": "bar"}
	wg.Add(1)
	ws.Dispatch("testev", payload)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for webhook dispatch")
	}

	received := <-payloadCh
	var got map[string]interface{}
	err = json.Unmarshal(received, &got)
	assert.NoError(t, err)
	assert.Equal(t, "bar", got["foo"])
}

func TestWebhookService_Dispatch_InvalidJSON(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svcIface := NewWebhookService(repos)
	ws := svcIface.(*webhookService)

	evs := map[model.EventType]bool{"ev": true}
	_, err := ws.Create("w2", "http://example.invalid", model.RequestTypePost, evs)
	assert.NoError(t, err)
	payload := make(chan int)
	ws.Dispatch("ev", payload)
	time.Sleep(10 * time.Millisecond)
}

func TestWebhookService_Dispatch_RequestError(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svcIface := NewWebhookService(repos)
	ws := svcIface.(*webhookService)

	evs := map[model.EventType]bool{"ev2": true}
	_, err := ws.Create("w3", "http://invalid.invalid", model.RequestTypePost, evs)
	assert.NoError(t, err)

	ws.httpClient.Timeout = 1 * time.Millisecond
	ws.Dispatch("ev2", map[string]string{"x": "y"})
	time.Sleep(50 * time.Millisecond)
}
