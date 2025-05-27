package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/service"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (h *Controller) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api",
		middleware.ApikeyAuth(h.services.Apikey),
	)
	api.GET("/collections/:alias/contents", h.listContents)
}

func (h *Controller) listContents(c *gin.Context) {
	alias := c.Param("alias")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	const perPage = 100
	offset := (page - 1) * perPage

	out, err := h.services.Api.ListByCollectionAlias(alias, offset, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"perPage": perPage,
		"data":    out,
	})
}
