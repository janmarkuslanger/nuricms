package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/handler"
	"github.com/janmarkuslanger/nuricms/internal/server"
)

type editDummy struct {
	ID   uint
	Name string
}

type mockEditService struct {
	Updated bool
	Err     error
}

func (m *mockEditService) UpdateByID(id uint, item editDummy) (*editDummy, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	m.Updated = true
	return &item, nil
}

func makeEditContext(method, path string) server.Context {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	return server.Context{Writer: w, Request: req}
}

func TestHandleEdit_Success(t *testing.T) {
	svc := &mockEditService{}
	ctx := makeEditContext("POST", "/")
	dto := editDummy{ID: 1, Name: "test"}
	handler.HandleEdit(ctx, svc, "1", dto, handler.HandlerOptions{
		RedirectOnSuccess: "/success",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)

	if resp.Code != http.StatusSeeOther || resp.Header().Get("Location") != "/success" {
		t.Errorf("expected redirect to /success, got: %d, %s", resp.Code, resp.Header().Get("Location"))
	}
	if !svc.Updated {
		t.Error("expected UpdateByID to be called")
	}
}

func TestHandleEdit_InvalidID_Render(t *testing.T) {
	svc := &mockEditService{}
	ctx := makeEditContext("POST", "/")
	handler.HandleEdit(ctx, svc, "abc", editDummy{}, handler.HandlerOptions{
		RenderOnFail: "fail.tmpl",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)

	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "fail.tmpl") {
		t.Errorf("expected render of fail.tmpl, got: %s", resp.Body.String())
	}
}

func TestHandleEdit_InvalidID_Redirect(t *testing.T) {
	svc := &mockEditService{}
	ctx := makeEditContext("POST", "/")
	handler.HandleEdit(ctx, svc, "abc", editDummy{}, handler.HandlerOptions{
		RedirectOnFail: "/fail-redirect",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)

	if resp.Code != http.StatusSeeOther || resp.Header().Get("Location") != "/fail-redirect" {
		t.Errorf("expected redirect to /fail-redirect, got: %d, %s", resp.Code, resp.Header().Get("Location"))
	}
}

func TestHandleEdit_UpdateError_Render(t *testing.T) {
	svc := &mockEditService{Err: errors.New("update failed")}
	ctx := makeEditContext("POST", "/")
	handler.HandleEdit(ctx, svc, "1", editDummy{}, handler.HandlerOptions{
		RenderOnFail: "updatefail.tmpl",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)

	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "updatefail.tmpl") {
		t.Errorf("expected render of updatefail.tmpl, got: %s", resp.Body.String())
	}
}

func TestHandleEdit_UpdateError_Redirect(t *testing.T) {
	svc := &mockEditService{Err: errors.New("update failed")}
	ctx := makeEditContext("POST", "/")
	handler.HandleEdit(ctx, svc, "1", editDummy{}, handler.HandlerOptions{
		RedirectOnFail: "/update-error",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)

	if resp.Code != http.StatusSeeOther || resp.Header().Get("Location") != "/update-error" {
		t.Errorf("expected redirect to /update-error, got: %d, %s", resp.Code, resp.Header().Get("Location"))
	}
}
