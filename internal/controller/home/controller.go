package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct Controller) RegisterRoutes(s *server.Server) {
	s.Handle("GET ", ct.home, middleware.Userauth(ct.services.User))
}

func (ct Controller) home(ctx server.Context) {
	utils.RenderWithLayoutHTTP(ctx, "home/home.tmpl", gin.H{}, http.StatusOK)
}
