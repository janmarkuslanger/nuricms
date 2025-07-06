package handler

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type CreateHandler[A any, B any] interface {
	Create(data A) (*B, error)
}

func HandleCreate[A any, B any](ctx server.Context, s CreateHandler[A, B], dto A, ho HandlerOptions) {
	if _, err := s.Create(dto); err != nil {
		if ho.RenderOnFail != "" {
			utils.RenderWithLayoutHTTP(ctx, ho.RenderOnFail, map[string]any{}, http.StatusOK)
			return
		}
	}

	if ho.RedirectOnSuccess != "" {
		http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnSuccess, http.StatusSeeOther)
	}
}
