package field

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
	secure := r.Group("/fields", middleware.Userauth(h.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.listFields)
	secure.GET("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showCreateField)
	secure.POST("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.createField)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showEditField)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.editField)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.deleteField)
}

func (h *Controller) listFields(c *gin.Context) {
	var fields []model.Field
	if err := db.DB.Preload("Collection").Find(&fields).Error; err != nil {
		utils.RenderWithLayout(c, "field/index.tmpl", gin.H{
			"error": "Failed to retrieve fields.",
		}, http.StatusInternalServerError)
		return
	}

	utils.RenderWithLayout(c, "field/index.tmpl", gin.H{
		"fields": fields,
	}, http.StatusOK)
}

func (h *Controller) showCreateField(c *gin.Context) {
	var collections []model.Collection
	if err := db.DB.Find(&collections).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	utils.RenderWithLayout(c, "field/create_or_edit.tmpl", gin.H{
		"collections": collections,
		"types": []model.FieldType{
			model.FieldTypeText, model.FieldTypeNumber, model.FieldTypeBoolean,
			model.FieldTypeDate, model.FieldTypeAsset, model.FieldTypeCollection,
			model.FieldTypeTextarea, model.FieldTypeRichText,
		},
	}, http.StatusOK)
}

func (h *Controller) showEditField(c *gin.Context) {
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

	utils.RenderWithLayout(c, "field/create_or_edit.tmpl", gin.H{
		"field":       field,
		"collections": collections,
		"types": []model.FieldType{
			model.FieldTypeText, model.FieldTypeNumber, model.FieldTypeBoolean,
			model.FieldTypeDate, model.FieldTypeAsset, model.FieldTypeCollection,
			model.FieldTypeTextarea, model.FieldTypeRichText,
		},
	}, http.StatusOK)
}

func (h *Controller) createField(c *gin.Context) {
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
		utils.RenderWithLayout(c, "field/create_or_edit.tmpl", gin.H{
			"error": "Failed to create field.",
			"field": field,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (h *Controller) editField(c *gin.Context) {
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
		utils.RenderWithLayout(c, "field/create_or_edit.tmpl", gin.H{
			"error": "Failed to update field.",
			"field": field,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (h *Controller) deleteField(c *gin.Context) {
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
