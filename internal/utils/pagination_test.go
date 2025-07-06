package utils_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/utils"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name     string
		query    url.Values
		wantPage int
		wantSize int
	}{
		{
			name:     "valid pagination params",
			query:    url.Values{"page": {"2"}, "pageSize": {"20"}},
			wantPage: 2, wantSize: 20,
		},
		{
			name:     "missing params",
			query:    url.Values{},
			wantPage: 1, wantSize: 10,
		},
		{
			name:     "invalid params",
			query:    url.Values{"page": {"abc"}, "pageSize": {"xyz"}},
			wantPage: 1, wantSize: 10,
		},
		{
			name:     "partial params",
			query:    url.Values{"page": {"3"}},
			wantPage: 3, wantSize: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					RawQuery: tt.query.Encode(),
				},
			}

			gotPage, gotSize := utils.ParsePagination(req)
			if gotPage != tt.wantPage || gotSize != tt.wantSize {
				t.Errorf("got (page: %d, size: %d), want (page: %d, size: %d)",
					gotPage, gotSize, tt.wantPage, tt.wantSize)
			}
		})
	}
}
