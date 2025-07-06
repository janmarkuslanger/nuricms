package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/handler"
	"github.com/janmarkuslanger/nuricms/internal/server"
)

type createDummy struct {
	Name string
}

type mockCreateService struct {
	ShouldFail bool
}

func (m *mockCreateService) Create(data createDummy) (*createDummy, error) {
	if m.ShouldFail {
		return nil, errors.New("create error")
	}
	return &data, nil
}

func makeCreateContext(method, path string) server.Context {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	return server.Context{Writer: w, Request: req}
}

func TestHandleCreate_Success(t *testing.T) {
	svc := &mockCreateService{}
	ctx := makeCreateContext("POST", "/create")

	handler.HandleCreate(ctx, svc, createDummy{Name: "ok"}, handler.HandlerOptions{
		RedirectOnSuccess: "/done",
	})

	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusSeeOther || resp.Header().Get("Location") != "/done" {
		t.Errorf("expected redirect to /done, got code %d and location %s", resp.Code, resp.Header().Get("Location"))
	}
}

func TestHandleCreate_Failure_Render(t *testing.T) {
	svc := &mockCreateService{ShouldFail: true}
	ctx := makeCreateContext("POST", "/create")

	handler.HandleCreate(ctx, svc, createDummy{Name: "fail"}, handler.HandlerOptions{
		RenderOnFail: "create_fail.tmpl",
	})

	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusOK || resp.Body.String() != "TEMPLATE: create_fail.tmpl" {
		t.Errorf("expected rendered template, got code %d and body %s", resp.Code, resp.Body.String())
	}
}

func TestHandleCreate_Failure_NoRender(t *testing.T) {
	svc := &mockCreateService{ShouldFail: true}
	ctx := makeCreateContext("POST", "/create")

	handler.HandleCreate(ctx, svc, createDummy{Name: "fail"}, handler.HandlerOptions{})

	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusOK && resp.Code != 0 {
		t.Errorf("expected default response code, got %d", resp.Code)
	}
}
