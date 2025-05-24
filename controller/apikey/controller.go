package apikey

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/model"
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
	secure := r.Group("/apikeys", middleware.Userauth(h.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleAdmin), h.showApikeys)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), h.showCreateApikey)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), h.createApikey)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), h.showEditApikey)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), h.deleteApikey)
}

func (h *Controller) showApikeys(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeNum = 10
	}

	keys, totalCount, err := h.services.Apikey.List(pageNum, pageSizeNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (totalCount + int64(pageSizeNum) - 1) / int64(pageSizeNum)

	utils.RenderWithLayout(c, "/apikey/index.tmpl", gin.H{
		"Apikeys":     keys,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": pageNum,
		"PageSize":    pageSizeNum,
	}, http.StatusOK)
}

func (h *Controller) showCreateApikey(c *gin.Context) {
	utils.RenderWithLayout(c, "/apikey/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (h *Controller) createApikey(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.Redirect(http.StatusSeeOther, "/apikeys")
	}
	h.services.Apikey.Create(name, 0)
	c.Redirect(http.StatusSeeOther, "/apikeys")
}

func (h *Controller) showEditApikey(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/apikeys")
	}

	apikey, err := h.services.Apikey.FindByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/apikeys")
	}

	utils.RenderWithLayout(c, "/apikey/create_or_edit.tmpl", gin.H{
		"Apikey": apikey,
	}, http.StatusOK)
}

func (h *Controller) deleteApikey(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/apikeys")
	}

	h.services.Apikey.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/apikeys")
}
