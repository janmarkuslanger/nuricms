package field

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/db"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
	"github.com/janmarkuslanger/nuricms/model"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/fields", middleware.Userauth(ct.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.listFields)
	secure.GET("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showCreateField)
	secure.POST("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.createField)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showEditField)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.editField)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.deleteField)
}

func (ct *Controller) listFields(c *gin.Context) {
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

	fields, totalCount, err := ct.services.Field.List(pageNum, pageSizeNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve fields."})
		return
	}

	totalPages := (totalCount + int64(pageSizeNum) - 1) / int64(pageSizeNum)

	utils.RenderWithLayout(c, "field/index.tmpl", gin.H{
		"Fields":      fields,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": pageNum,
		"PageSize":    pageSizeNum,
	}, http.StatusOK)
}

func (ct *Controller) showCreateField(c *gin.Context) {
	var collections []model.Collection
	if err := db.DB.Find(&collections).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	utils.RenderWithLayout(c, "field/create_or_edit.tmpl", gin.H{
		"Collections": collections,
		"Types": []model.FieldType{
			model.FieldTypeText, model.FieldTypeNumber, model.FieldTypeBoolean,
			model.FieldTypeDate, model.FieldTypeAsset, model.FieldTypeCollection,
			model.FieldTypeTextarea, model.FieldTypeRichText,
		},
	}, http.StatusOK)
}

func (ct *Controller) showEditField(c *gin.Context) {
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
		"Field":       field,
		"Collections": collections,
		"Types": []model.FieldType{
			model.FieldTypeText, model.FieldTypeNumber, model.FieldTypeBoolean,
			model.FieldTypeDate, model.FieldTypeAsset, model.FieldTypeCollection,
			model.FieldTypeTextarea, model.FieldTypeRichText,
		},
	}, http.StatusOK)
}

func (ct *Controller) createField(c *gin.Context) {
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
			"Error": "Failed to create field.",
			"Field": field,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (ct *Controller) editField(c *gin.Context) {
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
			"Error": "Failed to update field.",
			"Field": field,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (ct *Controller) deleteField(c *gin.Context) {
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
