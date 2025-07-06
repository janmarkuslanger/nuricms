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

type dummy struct {
	ID   uint
	Name string
}

type mockShowEditService struct {
	item *dummy
	err  error
}

func (m *mockShowEditService) FindByID(id uint) (*dummy, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.item, nil
}

func makeContext(method, path string) server.Context {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	return server.Context{Writer: w, Request: req}
}

func TestHandleShowEdit_Success(t *testing.T) {
	svc := &mockShowEditService{item: &dummy{ID: 1, Name: "Test"}}
	ctx := makeContext("GET", "/")
	handler.HandleShowEdit[dummy](ctx, svc, "1", handler.HandlerOptions{
		RenderOnSuccess: "success.tmpl",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "success.tmpl") {
		t.Errorf("expected success.tmpl, got: %s", resp.Body.String())
	}
}

func TestHandleShowEdit_InvalidID_Render(t *testing.T) {
	svc := &mockShowEditService{}
	ctx := makeContext("GET", "/")
	handler.HandleShowEdit[dummy](ctx, svc, "abc", handler.HandlerOptions{
		RenderOnFail: "fail.tmpl",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "fail.tmpl") {
		t.Errorf("expected fail.tmpl, got: %s", resp.Body.String())
	}
}

func TestHandleShowEdit_InvalidID_Redirect(t *testing.T) {
	svc := &mockShowEditService{}
	ctx := makeContext("GET", "/")
	handler.HandleShowEdit[dummy](ctx, svc, "abc", handler.HandlerOptions{
		RedirectOnFail: "/redirect",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusSeeOther || resp.Header().Get("Location") != "/redirect" {
		t.Errorf("expected redirect to /redirect, got: %d, %s", resp.Code, resp.Header().Get("Location"))
	}
}

func TestHandleShowEdit_FindError_Render(t *testing.T) {
	svc := &mockShowEditService{err: errors.New("not found")}
	ctx := makeContext("GET", "/")
	handler.HandleShowEdit[dummy](ctx, svc, "1", handler.HandlerOptions{
		RenderOnFail: "fail2.tmpl",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "fail2.tmpl") {
		t.Errorf("expected fail2.tmpl, got: %s", resp.Body.String())
	}
}

func TestHandleShowEdit_FindError_Redirect(t *testing.T) {
	svc := &mockShowEditService{err: errors.New("not found")}
	ctx := makeContext("GET", "/")
	handler.HandleShowEdit[dummy](ctx, svc, "1", handler.HandlerOptions{
		RedirectOnFail: "/not-found",
	})
	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusSeeOther || resp.Header().Get("Location") != "/not-found" {
		t.Errorf("expected redirect to /not-found, got: %d, %s", resp.Code, resp.Header().Get("Location"))
	}
}
