package utils

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
)

func GetParamOrRedirect(ctx server.Context, dest string, param string) (uint, bool) {
	v, ok := StringToUint(ctx.Request.PathValue(param))
	if !ok {
		http.Redirect(ctx.Writer, ctx.Request, dest, http.StatusSeeOther)
	}

	return v, ok
}
