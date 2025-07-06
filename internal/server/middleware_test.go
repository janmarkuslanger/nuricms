package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain(t *testing.T) {
	var callOrder []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw1-before")
			next.ServeHTTP(w, r)
			callOrder = append(callOrder, "mw1-after")
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw2-before")
			next.ServeHTTP(w, r)
			callOrder = append(callOrder, "mw2-after")
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(handler, mw1, mw2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	chained.ServeHTTP(w, req)

	expectedOrder := []string{
		"mw1-before",
		"mw2-before",
		"handler",
		"mw2-after",
		"mw1-after",
	}

	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("expected %d calls, got %d", len(expectedOrder), len(callOrder))
	}

	for i := range callOrder {
		if callOrder[i] != expectedOrder[i] {
			t.Errorf("at index %d: expected %q, got %q", i, expectedOrder[i], callOrder[i])
		}
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Result().StatusCode)
	}
}
