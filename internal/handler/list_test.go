package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/handler"
	"github.com/janmarkuslanger/nuricms/internal/server"
)

type listDummy struct {
	ID   int
	Name string
}

type mockListService struct {
	Items       []listDummy
	TotalCount  int64
	ReturnError bool
}

func (m *mockListService) List(page, pageSize int) ([]listDummy, int64, error) {
	return m.Items, m.TotalCount, nil
}

func makeListContext(url string) server.Context {
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	return server.Context{Writer: w, Request: req}
}

func TestHandleList(t *testing.T) {
	service := &mockListService{
		Items: []listDummy{
			{ID: 1, Name: "Foo"},
			{ID: 2, Name: "Bar"},
		},
		TotalCount: 5,
	}

	ctx := makeListContext("/items?page=1&pageSize=2")

	handler.HandleList(ctx, service, "my_template.tmpl")

	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if got := resp.Body.String(); got != "TEMPLATE: my_template.tmpl" {
		t.Fatalf("unexpected response body: %s", got)
	}
}
