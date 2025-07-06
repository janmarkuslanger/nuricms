package handler

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type ShowEditHandler[T any] interface {
	FindByID(id uint) (*T, error)
}

func HandleShowEdit[T any](ctx server.Context, s ShowEditHandler[T], idStr string, ho HandlerOptions) {
	data := make(map[string]any)
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

	item, err := s.FindByID(id)
	if err != nil {
		if ho.RenderOnFail != "" {
			data["Error"] = "Could not find item"
			utils.RenderWithLayoutHTTP(ctx, ho.RenderOnFail, data, http.StatusOK)
			return
		}

		if ho.RedirectOnFail != "" {
			http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnFail, http.StatusSeeOther)
			return
		}
	}

	if ho.TemplateData != nil {
		d, _ := ho.TemplateData()
		for k, v := range d {
			data[k] = v
		}
	}

	if ho.RenderOnSuccess != "" {
		data["Item"] = item
		utils.RenderWithLayoutHTTP(ctx, ho.RenderOnSuccess, data, http.StatusOK)
	}
}
