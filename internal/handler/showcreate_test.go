package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/handler"
	"github.com/janmarkuslanger/nuricms/internal/server"
)

func TestHandleShowCreate_RendersWithTemplateData(t *testing.T) {
	ctx := server.Context{
		Request: httptest.NewRequest(http.MethodGet, "/create", nil),
		Writer:  httptest.NewRecorder(),
	}

	handler.HandleShowCreate(ctx, handler.HandlerOptions{
		RenderOnSuccess: "create.tmpl",
		TemplateData: func() (map[string]any, error) {
			return map[string]any{"title": "Create Something"}, nil
		},
	})

	resp := ctx.Writer.(*httptest.ResponseRecorder)
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "create.tmpl") {
		t.Errorf("expected rendered output to contain 'create.tmpl', got: %s", resp.Body.String())
	}
}
