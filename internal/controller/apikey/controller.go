package apikey

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/apikeys", middleware.Userauth(ct.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleAdmin), ct.showApikeys)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), ct.showCreateApikey)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), ct.createApikey)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), ct.showEditApikey)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), ct.deleteApikey)
}

func (ct *Controller) showApikeys(c *gin.Context) {
	page, pageSize := utils.ParsePagination(c)

	keys, totalCount, err := ct.services.Apikey.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	utils.RenderWithLayout(c, "apikey/index.tmpl", gin.H{
		"Apikeys":     keys,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": page,
		"PageSize":    pageSize,
	}, http.StatusOK)
}

func (ct *Controller) showCreateApikey(c *gin.Context) {
	utils.RenderWithLayout(c, "apikey/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (ct *Controller) createApikey(c *gin.Context) {
	name := c.PostForm("name")
	_, err := ct.services.Apikey.CreateToken(name, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/apikeys/create")
		return
	}
	c.Redirect(http.StatusSeeOther, "/apikeys")
}

func (ct *Controller) showEditApikey(c *gin.Context) {
	id, ok := utils.GetParamOrRedirect(c, "/apikeys", "id")
	if !ok {
		return
	}

	apikey, err := ct.services.Apikey.FindByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/apikeys")
		return
	}

	utils.RenderWithLayout(c, "apikey/create_or_edit.tmpl", gin.H{
		"Apikey": apikey,
	}, http.StatusOK)
}

func (ct *Controller) deleteApikey(c *gin.Context) {
	id, ok := utils.GetParamOrRedirect(c, "/apikeys", "id")
	if !ok {
		return
	}

	ct.services.Apikey.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/apikeys")
}
