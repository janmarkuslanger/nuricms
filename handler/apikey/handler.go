package apikey

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/model"
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
	secure := r.Group("/apikeys", middleware.Userauth(h.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleAdmin), h.showApikeys)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), h.showCreateApikey)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), h.createApikey)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), h.deleteApikey)
}

func (h *Handler) showApikeys(c *gin.Context) {
	keys, _ := h.services.Apikey.List()
	utils.RenderWithLayout(c, "/apikey/index.tmpl", gin.H{
		"Apikeys": keys,
	}, http.StatusOK)
}

func (h *Handler) showCreateApikey(c *gin.Context) {
	utils.RenderWithLayout(c, "/apikey/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (h *Handler) createApikey(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.Redirect(http.StatusSeeOther, "/apikeys")
	}
	h.services.Apikey.Create(name, 0)
	c.Redirect(http.StatusSeeOther, "/apikeys")
}

func (h *Handler) deleteApikey(c *gin.Context) {

}
