package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/dto"
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

func (ct *Controller) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api",
		middleware.ApikeyAuth(ct.services.Apikey),
	)
	api.GET("/collections/:alias/content", ct.listContents)
	api.GET("/content/:id", ct.findContentById)
	api.GET("/collections/:alias/content/filter", ct.listContentsByFieldValue)
}

func (ct *Controller) findContentById(c *gin.Context) {
	id, _ := utils.StringToUint(c.Param("id"))
	data, _ := ct.services.Api.FindContentByID(id)
	c.JSON(http.StatusOK, dto.ApiResponse{
		Data:    data,
		Success: true,
		Meta: &dto.MetaData{
			Timestamp: time.Now().UTC(),
		},
	})
}

func (ct *Controller) listContents(c *gin.Context) {
	alias := c.Param("alias")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	const perPage = 100
	offset := (page - 1) * perPage

	data, err := ct.services.Api.FindContentByCollectionAlias(alias, offset, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Data:    data,
		Success: true,
		Meta: &dto.MetaData{
			Timestamp: time.Now().UTC(),
		},
		Pagination: &dto.Pagination{
			PerPage: perPage,
			Page:    page,
		},
	})
}

func (ct *Controller) listContentsByFieldValue(c *gin.Context) {
	alias := c.Param("alias")
	fieldAlias := c.Query("field")
	value := c.Query("value")
	if fieldAlias == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "field and value query parameters are required",
			"meta":    &dto.MetaData{Timestamp: time.Now().UTC()},
		})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	const perPage = 100
	offset := (page - 1) * perPage

	items, _ := ct.services.Api.FindContentByCollectionAndFieldValue(alias, fieldAlias, value, offset, perPage)

	c.JSON(http.StatusOK, dto.ApiResponse{
		Data:    items,
		Success: true,
		Meta: &dto.MetaData{
			Timestamp: time.Now().UTC(),
		},
		Pagination: &dto.Pagination{
			PerPage: perPage,
			Page:    page,
		},
	})
}
