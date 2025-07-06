package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/handler"
	"github.com/janmarkuslanger/nuricms/internal/server"
)

type mockDeleteHandler struct {
	DeleteCalled bool
	ShouldFail   bool
}

func (m *mockDeleteHandler) DeleteByID(id uint) error {
	m.DeleteCalled = true
	if m.ShouldFail {
		return errors.New("delete failed")
	}
	return nil
}

func TestHandleDelete_Success(t *testing.T) {
	mock := &mockDeleteHandler{}
	req := httptest.NewRequest(http.MethodPost, "/delete", nil)
	rr := httptest.NewRecorder()
	ctx := server.Context{Writer: rr, Request: req}

	handler.HandleDelete(ctx, mock, "42", handler.HandlerOptions{
		RedirectOnSuccess: "/success",
	})

	if !mock.DeleteCalled {
		t.Error("expected DeleteByID to be called")
	}

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rr.Code)
	}

	if loc := rr.Header().Get("Location"); loc != "/success" {
		t.Errorf("expected redirect to /success, got %s", loc)
	}
}

func TestHandleDelete_InvalidID(t *testing.T) {
	mock := &mockDeleteHandler{}
	req := httptest.NewRequest(http.MethodPost, "/delete", nil)
	rr := httptest.NewRecorder()
	ctx := server.Context{Writer: rr, Request: req}

	handler.HandleDelete(ctx, mock, "abc", handler.HandlerOptions{
		RedirectOnFail: "/fail",
	})

	if mock.DeleteCalled {
		t.Error("expected DeleteByID not to be called")
	}

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rr.Code)
	}

	if loc := rr.Header().Get("Location"); loc != "/fail" {
		t.Errorf("expected redirect to /fail, got %s", loc)
	}
}

func TestHandleDelete_DeleteFails(t *testing.T) {
	mock := &mockDeleteHandler{ShouldFail: true}
	req := httptest.NewRequest(http.MethodPost, "/delete", nil)
	rr := httptest.NewRecorder()
	ctx := server.Context{Writer: rr, Request: req}

	handler.HandleDelete(ctx, mock, "42", handler.HandlerOptions{
		RedirectOnFail: "/fail",
	})

	if !mock.DeleteCalled {
		t.Error("expected DeleteByID to be called")
	}

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rr.Code)
	}

	if loc := rr.Header().Get("Location"); loc != "/fail" {
		t.Errorf("expected redirect to /fail, got %s", loc)
	}
}
