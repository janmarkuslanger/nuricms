package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/server"
)

func TestHandleAndAddHandler(t *testing.T) {
	srv := server.NewServer()

	called := false
	srv.Handle("/test", func(ctx server.Context) {
		called = true
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write([]byte("OK"))
	})

	tests := []string{"/test", "/test/"}
	for _, path := range tests {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()

		srv.Mux.ServeHTTP(rr, req)

		resp := rr.Result()
		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		if string(body) != "OK" {
			t.Errorf("Expected body 'OK', got %q", string(body))
		}
	}

	if !called {
		t.Error("Handler was not called")
	}
}

func TestAddHandlerWithMiddleware(t *testing.T) {
	srv := server.NewServer()

	var steps []string

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			steps = append(steps, "middleware")
			next.ServeHTTP(w, r)
		})
	}

	handler := func(ctx server.Context) {
		steps = append(steps, "handler")
		ctx.Writer.WriteHeader(http.StatusOK)
	}

	srv.AddHandler("/mw", handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/mw", nil)
	rr := httptest.NewRecorder()
	srv.Mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rr.Code)
	}

	expected := []string{"middleware", "handler"}
	if strings.Join(steps, ",") != strings.Join(expected, ",") {
		t.Errorf("Expected call chain %v, got %v", expected, steps)
	}
}
