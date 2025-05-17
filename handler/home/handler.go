package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
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
	secure := r.Group("/", middleware.Userauth(h.services.User))
	secure.GET("", h.home)
}

func (h *Handler) home(c *gin.Context) {
	utils.RenderWithLayout(c, "home.tmpl", gin.H{}, http.StatusOK)
}
