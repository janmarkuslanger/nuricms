package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (h *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/", middleware.Userauth(h.services.User))
	secure.GET("", h.home)
}

func (h *Controller) home(c *gin.Context) {
	utils.RenderWithLayout(c, "home.tmpl", gin.H{}, http.StatusOK)
}
