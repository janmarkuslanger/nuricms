package client_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/janmarkuslanger/nuricms/pkg/client"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestNew_TrimsTrailingSlash(t *testing.T) {
	c := client.New("http://example.com/", "k")
	if got := c.BaseURL; got != "http://example.com" {
		t.Fatalf("expected trimmed base URL, got %q", got)
	}
}

func TestFindContentByID_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/content/42", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "secret" {
			t.Fatalf("missing or wrong API key header")
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{
			"success": true,
			"data": { "id": 42, "created_at": "now", "updated_at": "now", "values": {"title":"ok"} },
			"meta": { "timestamp": "2024-01-01T00:00:00Z" }
		}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "secret")
	item, err := api.FindContentByID(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != 42 {
		t.Fatalf("expected id 42, got %d", item.ID)
	}
}

func TestFindContentByID_Non200(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/content/7", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "k")
	_, err := api.FindContentByID(7)
	if err == nil || !strings.Contains(err.Error(), "API request failed") {
		t.Fatalf("expected non-200 error, got %v", err)
	}
}

func TestFindContentByID_BadJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/content/1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{")) // bad response
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "k")
	_, err := api.FindContentByID(1)
	if err == nil {
		t.Fatalf("expected JSON decode error")
	}
}

func TestFindContentByID_SuccessFalseOrNilData(t *testing.T) {
	mux := http.NewServeMux()
	// success=false
	mux.HandleFunc("/api/content/2", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success": false, "data": null}`)
	})
	// success=true, data=nil
	mux.HandleFunc("/api/content/3", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success": true, "data": null}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "k")
	if _, err := api.FindContentByID(2); err == nil {
		t.Fatalf("expected error for success=false")
	}
	if _, err := api.FindContentByID(3); err == nil {
		t.Fatalf("expected error for data=nil")
	}
}

func TestFindContentByCollectionAlias_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/collections/blog/content", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("page") != "2" || q.Get("perPage") != "5" {
			t.Fatalf("unexpected pagination: %v", q)
		}
		io.WriteString(w, `{
			"success": true,
			"data": [
				{"id":1,"created_at":"a","updated_at":"b","values":{"title":"t1"}},
				{"id":2,"created_at":"c","updated_at":"d","values":{"title":"t2"}}
			],
			"pagination":{"per_page":5,"page":2}
		}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "k")
	items, pag, err := api.FindContentByCollectionAlias("blog", 2, 5)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if pag == nil || pag.Page != 2 || pag.PerPage != 5 {
		t.Fatalf("unexpected pagination: %#v", pag)
	}
}

func TestFindContentByCollectionAlias_AliasEscaping(t *testing.T) {
	alias := "my blog/2025"
	path := "/api/collections/" + url.PathEscape(alias) + "/content"
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success": true, "data": [], "pagination":{"per_page":1,"page":1}}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "k")
	_, _, err := api.FindContentByCollectionAlias(alias, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFindContentByCollectionAndFieldValue_Success(t *testing.T) {
	field := "slug"
	value := "mein artikel?"
	alias := "blog"

	expectedPath := "/api/collections/" + url.PathEscape(alias) + "/content/filter"

	mux := http.NewServeMux()
	mux.HandleFunc(expectedPath, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("field") != field || q.Get("value") != value || q.Get("page") != "1" || q.Get("perPage") != "10" {
			t.Fatalf("unexpected query: %v", q)
		}
		io.WriteString(w, `{
			"success": true,
			"data": [{"id":99,"created_at":"x","updated_at":"y","values":{"slug":"mein artikel?"}}],
			"pagination":{"per_page":10,"page":1}
		}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "k")
	items, pag, err := api.FindContentByCollectionAndFieldValue(alias, field, value, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].ID != 99 {
		t.Fatalf("unexpected items: %#v", items)
	}
	if pag == nil || pag.Page != 1 || pag.PerPage != 10 {
		t.Fatalf("unexpected pagination: %#v", pag)
	}
}

func TestGet_HTTPClientError(t *testing.T) {
	api := client.New("http://example.com", "k")
	api.HTTPClient = &http.Client{
		Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return nil, errors.New("boom")
		}),
		Timeout: 1 * time.Second,
	}
	_, err := api.FindContentByID(1)
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGet_JSONDecodeErrorOnList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/collections/blog/content", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"success": true, "data": [`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "k")
	_, _, err := api.FindContentByCollectionAlias("blog", 1, 1)
	if err == nil {
		t.Fatalf("expected JSON decode error")
	}
}

func TestHeadersAreSet(t *testing.T) {
	var gotAPIKey, gotAccept string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/content/5", func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("X-API-Key")
		gotAccept = r.Header.Get("Accept")
		io.WriteString(w, `{"success": true, "data":{"id":5,"created_at":"","updated_at":"","values":{}}}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := client.New(srv.URL, "superkey")
	if _, err := api.FindContentByID(5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAPIKey != "superkey" || gotAccept != "application/json" {
		t.Fatalf("unexpected headers: X-API-Key=%q Accept=%q", gotAPIKey, gotAccept)
	}
}
