package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderTemplate_Success(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "hello.tmpl")
	err := os.WriteFile(tpl, []byte("Hello {{.Name}}!"), 0644)
	require.NoError(t, err)

	out, err := RenderTemplate(tpl, map[string]string{"Name": "World"})
	require.NoError(t, err)
	require.Equal(t, "Hello World!", out)
}

func TestRenderTemplate_FileNotFound(t *testing.T) {
	_, err := RenderTemplate("doesnotexist.tmpl", nil)
	require.Error(t, err)
}

func TestRenderTemplate_ParseError(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "bad.tmpl")
	err := os.WriteFile(tpl, []byte("{{"), 0644)
	require.NoError(t, err)

	_, err = RenderTemplate(tpl, nil)
	require.Error(t, err)
}
