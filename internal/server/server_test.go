package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/stretchr/testify/require"
)

func TestServer_Handle(t *testing.T) {
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

func TestServer_Handle_Middleware(t *testing.T) {
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

func TestServer_ServeHTTP_Success(t *testing.T) {
	srv := server.NewServer()

	srv.Handle("/test", func(ctx server.Context) {
		ctx.Writer.WriteHeader(http.StatusTeapot)
		ctx.Writer.Write([]byte("hello world"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusTeapot, rec.Code)
	require.Equal(t, "hello world", rec.Body.String())
}

func TestServer_Static_GET(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "nested"), 0o755))
	want := []byte("hello")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "nested", "file.txt"), want, 0o644))

	srv := server.NewServer()
	srv.Static("/public/assets", dir)

	req := httptest.NewRequest(http.MethodGet, "/public/assets/nested/file.txt", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	body, _ := io.ReadAll(rec.Body)
	require.Equal(t, string(want), string(body))
}

func TestServer_Static_HEAD(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0o644))

	srv := server.NewServer()
	srv.Static("/public/assets", dir)

	req := httptest.NewRequest(http.MethodHead, "/public/assets/file.txt", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, 0, rec.Body.Len())
	require.NotEmpty(t, rec.Header().Get("Content-Length"))
}

func TestServer_Static_MethodNotAllowedLike(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("x"), 0o644))

	srv := server.NewServer()
	srv.Static("/public/assets", dir)

	req := httptest.NewRequest(http.MethodPost, "/public/assets/file.txt", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
