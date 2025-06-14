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
	svc := NewWebhookService(repos)

	events := map[model.EventType]bool{
		model.EventType("ev1"): true,
		model.EventType("ev2"): false,
		model.EventType("ev3"): true,
	}
	wh, err := svc.Create("name1", "http://example.com", model.RequestType("POST"), events)
	assert.NoError(t, err)
	assert.NotZero(t, wh.ID)
	assert.Contains(t, wh.Events, "ev1,")
	assert.Contains(t, wh.Events, "ev3,")
	assert.NotContains(t, wh.Events, "ev2,")

	wh2, err := svc.Create("name2", "http://example.org", model.RequestType("GET"), map[model.EventType]bool{"e": true})
	assert.NoError(t, err)
	list, total, err := svc.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	ids := []uint{list[0].ID, list[1].ID}
	assert.Contains(t, ids, wh.ID)
	assert.Contains(t, ids, wh2.ID)

	found, err := svc.FindByID(wh.ID)
	assert.NoError(t, err)
	assert.Equal(t, wh.ID, found.ID)

	found.Url = "http://changed"
	err = svc.Save(found)
	assert.NoError(t, err)
	reloaded, err := svc.FindByID(wh.ID)
	assert.NoError(t, err)
	assert.Equal(t, "http://changed", reloaded.Url)

	err = svc.DeleteByID(wh2.ID)
	assert.NoError(t, err)
	_, err = repos.Webhook.FindByID(wh2.ID)
	assert.Error(t, err)
}

func TestWebhookService_Dispatch_NoHooks(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewWebhookService(repos)
	svc.Dispatch("noevent", map[string]string{"k": "v"})
	time.Sleep(10 * time.Millisecond)
}

func TestWebhookService_Dispatch_SendsRequest(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewWebhookService(repos)
	var wg sync.WaitGroup

	payloadCh := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		body, _ := io.ReadAll(r.Body)
		payloadCh <- body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc.httpClient = server.Client()

	evs := map[model.EventType]bool{"testev": true}
	_, err := svc.Create("w1", server.URL, model.RequestType("POST"), evs)
	assert.NoError(t, err)

	payload := map[string]interface{}{"foo": "bar"}
	wg.Add(1)
	svc.Dispatch("testev", payload)

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
	svc := NewWebhookService(repos)

	evs := map[model.EventType]bool{"ev": true}
	_, err := svc.Create("w2", "http://invalid.local", model.RequestType("POST"), evs)
	assert.NoError(t, err)

	payload := make(chan int)
	svc.Dispatch("ev", payload)
	time.Sleep(10 * time.Millisecond)
}

func TestWebhookService_Dispatch_RequestError(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewWebhookService(repos)

	evs := map[model.EventType]bool{"ev2": true}
	_, err := svc.Create("w3", "http://invalid.invalid", model.RequestType("POST"), evs)
	assert.NoError(t, err)
	svc.httpClient.Timeout = 10 * time.Millisecond
	svc.Dispatch("ev2", map[string]string{"x": "y"})
	time.Sleep(50 * time.Millisecond)
}
