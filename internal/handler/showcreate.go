package handler

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

func HandleShowCreate(ctx server.Context, ho HandlerOptions) {
	data := make(map[string]any)

	if ho.TemplateData != nil {
		d, _ := ho.TemplateData()
		for k, v := range d {
			data[k] = v
		}
	}

	if ho.RenderOnSuccess != "" {
		utils.RenderWithLayoutHTTP(ctx, ho.RenderOnSuccess, data, http.StatusOK)
	}

}
