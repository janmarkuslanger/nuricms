package handler

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type EditHandler[A any, B any] interface {
	UpdateByID(id uint, item A) (*B, error)
}

func HandleEdit[A any, B any](ctx server.Context, s EditHandler[A, B], idStr string, dto A, ho HandlerOptions) {
	id, ok := utils.StringToUint(idStr)
	if !ok {

		if ho.RenderOnFail != "" {
			utils.RenderWithLayoutHTTP(ctx, ho.RenderOnFail, map[string]any{
				"Error": "ID wrong or could not be converted",
			}, http.StatusOK)
			return
		}

		if ho.RedirectOnFail != "" {
			http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnFail, http.StatusSeeOther)
			return
		}
	}

	_, err := s.UpdateByID(id, dto)
	if err != nil {
		if ho.RenderOnFail != "" {
			utils.RenderWithLayoutHTTP(ctx, ho.RenderOnFail, map[string]any{
				"Error": "Could not find item",
			}, http.StatusOK)
			return
		}

		if ho.RedirectOnFail != "" {
			http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnFail, http.StatusSeeOther)
			return
		}
	}

	if ho.RedirectOnSuccess != "" {
		http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnSuccess, http.StatusSeeOther)
	}
}
