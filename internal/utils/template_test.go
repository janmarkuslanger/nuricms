package utils_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRenderWithLayoutHTTP_Success(t *testing.T) {
	fakeFS := fstest.MapFS{
		"templates/base/layout.tmpl": &fstest.MapFile{
			Data: []byte(`{{ define "layout.tmpl" }}<html>{{ template "content" . }}</html>{{ end }}`),
		},
		"templates/test.tmpl": &fstest.MapFile{
			Data: []byte(`{{ define "content" }}Hello, {{ if .IsLoggedIn }}User{{ else }}Guest{{ end }}{{ end }}`),
		},
	}

	utils.SetTemplatesFS(fakeFS)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := server.Context{
		Request: req,
		Writer:  w,
	}

	utils.RenderWithLayoutHTTP(ctx, "test.tmpl", map[string]any{}, http.StatusOK)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf := new(strings.Builder)
	_, err := io.Copy(buf, res.Body)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "<html>Hello, Guest</html>")
}

func TestRenderWithLayoutHTTP_TemplateParseError(t *testing.T) {
	fakeFS := fstest.MapFS{
		"templates/base/layout.tmpl": &fstest.MapFile{
			Data: []byte(`{{ define "layout.tmpl" }}<html>{{ template "oops" . }}</html>{{ end }}`),
		},
		"templates/test.tmpl": &fstest.MapFile{
			Data: []byte(`{{ define "content" }}Hello{{ end }}`),
		},
	}
	utils.SetTemplatesFS(fakeFS)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := server.Context{
		Request: req,
		Writer:  w,
	}

	utils.RenderWithLayoutHTTP(ctx, "test.tmpl", map[string]any{}, http.StatusOK)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
