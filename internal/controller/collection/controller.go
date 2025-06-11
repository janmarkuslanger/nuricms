package collection

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/db"
	"github.com/janmarkuslanger/nuricms/internal/dto"
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
	secure := r.Group("/collections", middleware.Userauth(ct.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleAdmin, model.RoleEditor), ct.showCollections)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), ct.showCreateCollection)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), ct.createCollection)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), ct.showEditCollection)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleAdmin), ct.editCollection)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), ct.deleteCollection)
}

func (ct *Controller) showCollections(c *gin.Context) {
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

	collections, totalCount, err := ct.services.Collection.List(pageNum, pageSizeNum)

	totalPages := (totalCount + int64(pageSizeNum) - 1) / int64(pageSizeNum)

	utils.RenderWithLayout(c, "collection/index.tmpl", gin.H{
		"Collections": collections,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": pageNum,
		"PageSize":    pageSizeNum,
	}, http.StatusOK)
}

func (ct *Controller) showCreateCollection(c *gin.Context) {
	utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (ct *Controller) createCollection(c *gin.Context) {
	name := c.PostForm("name")
	alias := c.PostForm("alias")
	description := c.PostForm("description")

	if name == "" || alias == "" {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"error": "Name and Alias are required fields.",
		}, http.StatusOK)
		return
	}

	data := &dto.CollectionData{
		Name:        name,
		Alias:       alias,
		Description: description,
	}

	_, err := ct.services.Collection.Create(data)
	if err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"Error": "Failed to create collection.",
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/collections")
}

func (ct *Controller) showEditCollection(c *gin.Context) {
	idStr := c.Param("id")

	id, ok := utils.StringToUint(idStr)
	if !ok {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"Error": "Collection not found.",
		}, http.StatusInternalServerError)
		return
	}

	collection, err := ct.services.Collection.FindByID(id)
	if err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"Error": "Collection not found.",
		}, http.StatusInternalServerError)
		return
	}

	utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
		"Collection": collection,
	}, http.StatusOK)
}

func (ct *Controller) editCollection(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"Error": "Collection not found.",
		}, http.StatusNotFound)
		return
	}

	_, err := ct.services.Collection.UpdateByID(id, dto.CollectionData{
		Name:        c.PostForm("name"),
		Alias:       c.PostForm("alias"),
		Description: c.PostForm("description"),
	})
	if err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.tmpl", gin.H{
			"Error": "Collection could not be updated.",
		}, http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusSeeOther, "/collections")
}

func (ct *Controller) deleteCollection(c *gin.Context) {
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
