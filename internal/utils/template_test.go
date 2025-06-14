package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/embedfs"
	"github.com/stretchr/testify/assert"
)

func TestRenderWithLayout_HasUserID_FieldTrue(t *testing.T) {
	stubFS := fstest.MapFS{
		"templates/base/layout.tmpl": {Data: []byte(`{{if .IsLoggedIn}}IN{{else}}OUT{{end}}|{{template "content.tmpl" .}}`)},
		"templates/content.tmpl":     {Data: []byte(`{{define "content.tmpl"}}MSG: {{.Msg}}{{end}}`)},
	}

	SetTemplatesFS(stubFS)
	defer SetTemplatesFS(embedfs.TemplatesFS)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Set("userID", 42)

	RenderWithLayout(ctx, "content.tmpl", gin.H{"Msg": "Hello"}, http.StatusTeapot)

	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.Equal(t, "IN|MSG: Hello", rec.Body.String())
}

func TestRenderWithLayout_NoUserID_FieldFalse(t *testing.T) {
	stubFS := fstest.MapFS{
		"templates/base/layout.tmpl": {Data: []byte(`{{if .IsLoggedIn}}IN{{else}}OUT{{end}}|{{template "content.tmpl" .}}`)},
		"templates/content.tmpl":     {Data: []byte(`{{define "content.tmpl"}}MSG: {{.Msg}}{{end}}`)},
	}
	originalFS := embedfs.TemplatesFS
	SetTemplatesFS(stubFS)
	defer SetTemplatesFS(originalFS)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	RenderWithLayout(ctx, "content.tmpl", gin.H{"Msg": "World"}, http.StatusOK)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OUT|MSG: World", rec.Body.String())
}

func TestRenderWithLayout_TemplateExecutionError(t *testing.T) {
	stubFS := fstest.MapFS{
		"templates/base/layout.tmpl": {Data: []byte(`{{template "missing" .}}`)},
	}
	SetTemplatesFS(stubFS)
	defer SetTemplatesFS(embedfs.TemplatesFS)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	RenderWithLayout(ctx, "content.tmpl", gin.H{}, http.StatusOK)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Template parse error:")
}
