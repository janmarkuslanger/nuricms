package template_test

import (
	"embed"
	"errors"
	"io/fs"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/template"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var testFS embed.FS

func TestRenderTemplate_Success(t *testing.T) {
	result, err := template.RenderTemplate(testFS, "testdata/hello.tmpl", map[string]string{"Name": "World"})
	require.NoError(t, err)
	require.Equal(t, "Hello, World!", result)
}

func TestRenderTemplate_ParseErr(t *testing.T) {
	_, err := template.RenderTemplate(testFS, "testdata/false.tmpl", nil)
	require.Error(t, err)
}

func TestRenderTemplate_TemplateNotFound(t *testing.T) {
	_, err := template.RenderTemplate(testFS, "testdata/notfound.tmpl", nil)
	require.Error(t, err)
	require.True(t, errors.Is(err, fs.ErrNotExist) || err.Error() != "")
}
