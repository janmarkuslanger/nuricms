package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
)

type Handler struct {
	services *service.Set
}

func NewHandler(services *service.Set) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/assets/create", h.showCreateAsset)
}

func (h *Handler) showCreateAsset(c *gin.Context) {
	utils.RenderWithLayout(c, "/asset/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}
