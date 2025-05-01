package field

import (
	"net/http"
	"strconv"

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
	r.GET("/fields", h.listFields)
	r.GET("/fields/create", h.showCreateField)
	r.POST("/fields/create", h.createField)
	r.GET("/fields/edit/:id", h.showEditField)
	r.POST("/fields/edit/:id", h.editField)
	r.POST("/fields/delete/:id", h.deleteField)
}

func (h *Handler) listFields(c *gin.Context) {
	var fields []model.Field
	if err := db.DB.Preload("Collection").Find(&fields).Error; err != nil {
		utils.RenderWithLayout(c, "field/index.html", gin.H{
			"error": "Failed to retrieve fields.",
		}, http.StatusInternalServerError)
		return
	}

	utils.RenderWithLayout(c, "field/index.html", gin.H{
		"fields": fields,
	}, http.StatusOK)
}

func (h *Handler) showCreateField(c *gin.Context) {
	var collections []model.Collection
	if err := db.DB.Find(&collections).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	utils.RenderWithLayout(c, "field/create_or_edit.html", gin.H{
		"collections": collections,
		"types": []model.FieldType{
			model.FieldTypeText, model.FieldTypeNumber, model.FieldTypeBoolean,
			model.FieldTypeDate, model.FieldTypeAsset, model.FieldTypeCollection,
			model.FieldTypeTextarea, model.FieldTypeRichText,
		},
	}, http.StatusOK)
}

func (h *Handler) showEditField(c *gin.Context) {
	var field model.Field
	if err := db.DB.Preload("Collection").First(&field, c.Param("id")).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	var collections []model.Collection
	if err := db.DB.Find(&collections).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	utils.RenderWithLayout(c, "field/create_or_edit.html", gin.H{
		"field":       field,
		"collections": collections,
		"types": []model.FieldType{
			model.FieldTypeText, model.FieldTypeNumber, model.FieldTypeBoolean,
			model.FieldTypeDate, model.FieldTypeAsset, model.FieldTypeCollection,
			model.FieldTypeTextarea, model.FieldTypeRichText,
		},
	}, http.StatusOK)
}

func (h *Handler) createField(c *gin.Context) {
	name := c.PostForm("name")
	alias := c.PostForm("alias")
	collectionIDStr := c.PostForm("collection_id")

	if name == "" || alias == "" || collectionIDStr == "" {
		c.Redirect(http.StatusSeeOther, "/fields/")
		return
	}

	collectionID, _ := strconv.ParseUint(collectionIDStr, 10, 32)

	field := model.Field{
		Name:         name,
		Alias:        alias,
		CollectionID: uint(collectionID),
		FieldType:    model.FieldType(c.PostForm("field_type")),
		IsList:       c.PostForm("is_list") == "on",
		IsRequired:   c.PostForm("is_required") == "on",
		DisplayField: c.PostForm("display_field") == "on",
	}

	if err := db.DB.Create(&field).Error; err != nil {
		utils.RenderWithLayout(c, "field/create_or_edit.html", gin.H{
			"error": "Failed to create field.",
			"field": field,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (h *Handler) editField(c *gin.Context) {
	id := c.Param("id")
	var field model.Field
	if err := db.DB.First(&field, id).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	field.Name = c.PostForm("name")
	field.Alias = c.PostForm("alias")
	collectionID, _ := strconv.ParseUint(c.PostForm("collection_id"), 10, 32)
	field.CollectionID = uint(collectionID)
	field.FieldType = model.FieldType(c.PostForm("field_type"))
	field.IsList = c.PostForm("is_list") == "on"
	field.IsRequired = c.PostForm("is_required") == "on"
	field.DisplayField = c.PostForm("display_field") == "on"

	if err := db.DB.Save(&field).Error; err != nil {
		utils.RenderWithLayout(c, "field/create_or_edit.html", gin.H{
			"error": "Failed to update field.",
			"field": field,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (h *Handler) deleteField(c *gin.Context) {
	id := c.Param("id")

	var field model.Field
	if err := db.DB.First(&field, id).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	if err := db.DB.Delete(&field).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}
