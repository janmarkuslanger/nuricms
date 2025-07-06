package utils

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
)

func GetParamOrRedirect(ctx server.Context, r string, p string) (uint, bool) {
	id, ok := StringToUint(ctx.Request.PathValue(p))
	if !ok {
		http.Redirect(ctx.Writer, ctx.Request, r, http.StatusSeeOther)
	}

	return id, ok
}
