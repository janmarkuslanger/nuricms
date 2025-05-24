package collection

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/db"
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
	secure := r.Group("/collections", middleware.Userauth(h.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleAdmin, model.RoleEditor), h.showCollections)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), h.showCreateCollection)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), h.createCollection)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), h.showEditCollection)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleAdmin), h.editCollection)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), h.deleteCollection)
}

func (h *Controller) showCollections(c *gin.Context) {
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

	collections, totalCount, err := h.services.Collection.List(pageNum, pageSizeNum)

	totalPages := (totalCount + int64(pageSizeNum) - 1) / int64(pageSizeNum)

	utils.RenderWithLayout(c, "collection/index.tmpl", gin.H{
		"Collections": collections,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": pageNum,
		"PageSize":    pageSizeNum,
	}, http.StatusOK)
}

func (h *Controller) showCreateCollection(c *gin.Context) {
	utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (h *Controller) createCollection(c *gin.Context) {
	name := c.PostForm("name")
	alias := c.PostForm("alias")
	description := c.PostForm("description")

	if name == "" || alias == "" {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"error": "Name and Alias are required fields.",
		}, http.StatusOK)
		return
	}

	collection := model.Collection{
		Name:        name,
		Alias:       alias,
		Description: description,
	}

	if err := db.DB.Create(&collection).Error; err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"error": "Failed to create collection.",
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/collections")
}

func (h *Controller) showEditCollection(c *gin.Context) {
	id := c.Param("id")
	var collection model.Collection
	if err := db.DB.First(&collection, id).Error; err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"error": "Collection not found.",
		}, http.StatusInternalServerError)
		return
	}

	utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
		"Collection": collection,
	}, http.StatusOK)
}

func (h *Controller) editCollection(c *gin.Context) {
	id := c.Param("id")
	var collection model.Collection
	if err := db.DB.First(&collection, id).Error; err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"error": "Collection not found.",
		}, http.StatusNotFound)
		return
	}

	name := c.PostForm("name")
	alias := c.PostForm("alias")
	description := c.PostForm("description")

	if name == "" || alias == "" {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"error":      "Name and Alias are required fields.",
			"Collection": collection,
		}, http.StatusBadRequest)
		return
	}

	collection.Name = name
	collection.Alias = alias
	collection.Description = description

	if err := db.DB.Save(&collection).Error; err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"error":      "Failed to update.",
			"Collection": collection,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/collections")
}

func (h *Controller) deleteCollection(c *gin.Context) {
	id := c.Param("id")
	var collection model.Collection
	if err := db.DB.First(&collection, id).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	if err := db.DB.Delete(&collection).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	c.Redirect(http.StatusSeeOther, "/collections")
}
