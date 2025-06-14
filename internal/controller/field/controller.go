package field

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	collections, _, err := ct.services.Collection.List(1, 999999999999999999)
	if err != nil {
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
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	field, err := ct.services.Field.FindByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	collections, _, err := ct.services.Collection.List(1, 999999999999999999)
	if err != nil {
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
	_, err := ct.services.Field.Create(dto.FieldData{
		Name:         c.PostForm("name"),
		Alias:        c.PostForm("alias"),
		CollectionID: c.PostForm("collection_id"),
		FieldType:    c.PostForm("field_type"),
		IsList:       c.PostForm("is_list"),
		IsRequired:   c.PostForm("is_required"),
		DisplayField: c.PostForm("display_field"),
	})
	if err != nil {
		utils.RenderWithLayout(c, "field/create_or_edit.tmpl", gin.H{
			"Error": "Failed to create field.",
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (ct *Controller) editField(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	field, err := ct.services.Field.UpdateByID(id, dto.FieldData{
		Name:         c.PostForm("name"),
		Alias:        c.PostForm("alias"),
		CollectionID: c.PostForm("collection_id"),
		FieldType:    c.PostForm("field_type"),
		IsList:       c.PostForm("is_list"),
		IsRequired:   c.PostForm("is_required"),
		DisplayField: c.PostForm("display_field"),
	})
	if err != nil {
		utils.RenderWithLayout(c, "field/create_or_edit.tmpl", gin.H{
			"Error": "Failed to update field.",
			"Field": field,
		}, http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, "/fields")
}

func (ct *Controller) deleteField(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/fields")
		return
	}

	ct.services.Field.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/fields")
}
