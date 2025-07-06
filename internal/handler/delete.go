package handler

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type DeleteHandler interface {
	DeleteByID(id uint) error
}

func HandleDelete(ctx server.Context, s DeleteHandler, idStr string, ho HandlerOptions) {
	id, ok := utils.StringToUint(idStr)
	if !ok {
		if ho.RedirectOnFail != "" {
			http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnFail, http.StatusSeeOther)
			return
		}
	}

	if err := s.DeleteByID(id); err != nil {
		if ho.RedirectOnFail != "" {
			http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnFail, http.StatusSeeOther)
			return
		}
	}

	if ho.RedirectOnSuccess != "" {
		http.Redirect(ctx.Writer, ctx.Request, ho.RedirectOnSuccess, http.StatusSeeOther)
	}
}
