package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/db"
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
	r.GET("/collections", h.showCollections)
	r.GET("/collections/create", h.showCreateCollection)
	r.POST("/collections/create", h.createCollection)
	r.GET("/collections/edit/:id", h.showEditCollection)
	r.POST("/collections/edit/:id", h.editCollection)
	r.POST("/collections/delete/:id", h.deleteCollection)
}

func (h *Handler) showCollections(c *gin.Context) {
	var collections []model.Collection
	if err := db.DB.Find(&collections).Error; err != nil {
		utils.RenderWithLayout(c, "collection/index.html", gin.H{
			"error": "Failed to retrieve collections.",
		}, http.StatusInternalServerError)
		return
	}

	utils.RenderWithLayout(c, "collection/index.html", gin.H{
		"collections": collections,
	}, http.StatusOK)
}

func (h *Handler) showCreateCollection(c *gin.Context) {
	utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{}, http.StatusOK)
}

func (h *Handler) createCollection(c *gin.Context) {
	name := c.PostForm("name")
	alias := c.PostForm("alias")
	description := c.PostForm("description")

	if name == "" || alias == "" {
		utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{
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
		utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{
			"error": "Failed to create collection.",
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/collections")
}

func (h *Handler) showEditCollection(c *gin.Context) {
	id := c.Param("id")
	var collection model.Collection
	if err := db.DB.First(&collection, id).Error; err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{
			"error": "Collection not found.",
		}, http.StatusInternalServerError)
		return
	}

	utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{
		"collection": collection,
	}, http.StatusOK)
}

func (h *Handler) editCollection(c *gin.Context) {
	id := c.Param("id")
	var collection model.Collection
	if err := db.DB.First(&collection, id).Error; err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{
			"error": "Collection not found.",
		}, http.StatusNotFound)
		return
	}

	name := c.PostForm("name")
	alias := c.PostForm("alias")
	description := c.PostForm("description")

	if name == "" || alias == "" {
		utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{
			"error":      "Name and Alias are required fields.",
			"collection": collection,
		}, http.StatusBadRequest)
		return
	}

	collection.Name = name
	collection.Alias = alias
	collection.Description = description

	if err := db.DB.Save(&collection).Error; err != nil {
		utils.RenderWithLayout(c, "collection/create_or_edit.html", gin.H{
			"error":      "Failed to update.",
			"collection": collection,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/collections")
}

func (h *Handler) deleteCollection(c *gin.Context) {
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
