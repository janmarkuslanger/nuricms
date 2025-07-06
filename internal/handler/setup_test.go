package handler_test

import (
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

func init() {
	utils.RenderWithLayoutHTTP = func(ctx server.Context, tmpl string, data map[string]any, code int) {
		ctx.Writer.WriteHeader(code)
		_, _ = ctx.Writer.Write([]byte("TEMPLATE: " + tmpl))
	}
}
