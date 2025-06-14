package content

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
	secure := r.Group("/content/collections", middleware.Userauth(ct.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showCollections)
	secure.GET("/:id/show", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.listContent)
	secure.GET("/:id/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showCreateContent)
	secure.POST("/:id/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.createContent)
	secure.GET("/:id/edit/:contentID", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showEditContent)
	secure.POST("/:id/edit/:contentID", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.editContent)
	secure.POST("/:id/delete/:contentID", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.deleteContent)
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

	utils.RenderWithLayout(c, "content/collections.tmpl", gin.H{
		"Collections": collections,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": pageNum,
		"PageSize":    pageSizeNum,
	}, http.StatusOK)
}

func (ct *Controller) showCreateContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	fields, err := ct.services.Field.FindByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	collection, err := ct.services.Collection.FindByID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	contents, err := ct.services.Content.FindContentsWithDisplayContentValue()
	assets, _, err := ct.services.Asset.List(1, 100000)

	fieldsContent := make([]FieldContent, 0)
	for _, field := range fields {
		fieldsContent = append(fieldsContent, FieldContent{
			Field:   field,
			Content: contents,
			Assets:  assets,
		})
	}

	utils.RenderWithLayout(c, "content/create_or_edit.tmpl", gin.H{
		"FieldsHtml": RenderFields(fieldsContent),
		"Collection": collection,
	}, http.StatusOK)
}

func (ct *Controller) createContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if err := c.Request.ParseForm(); err != nil || !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	if _, err := ct.services.Content.CreateWithValues(dto.ContentWithValues{
		CollectionID: collectionID,
		FormData:     c.Request.PostForm,
	}); err == nil {
		ct.services.Webhook.Dispatch(string(model.EventContentCreated), nil)
	}

	c.Redirect(http.StatusSeeOther, "/content/collections")
}

func (ct *Controller) editContent(c *gin.Context) {
	colID, okCol := utils.StringToUint(c.Param("id"))
	conID, okCon := utils.StringToUint(c.Param("contentID"))
	err := c.Request.ParseForm()
	if !okCon || !okCol || err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	if _, err := ct.services.Content.EditWithValues(dto.ContentWithValues{
		CollectionID: colID,
		ContentID:    conID,
		FormData:     c.Request.PostForm,
	}); err != nil {
		ct.services.Webhook.Dispatch(string(model.EventContentUpdated), nil)
	}

	c.Redirect(http.StatusSeeOther, "/content/collections")
}

func (ct *Controller) listContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

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

	contents, totalCount, err := ct.services.Content.FindDisplayValueByCollectionID(collectionID, pageNum, pageSizeNum)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	groups := ContentsToContentGroup(contents)

	fields, err := ct.services.Field.FindDisplayFieldsByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	totalPages := (totalCount + int64(pageSizeNum) - 1) / int64(pageSizeNum)

	utils.RenderWithLayout(c, "content/content_list.tmpl", gin.H{
		"Groups":       groups,
		"Fields":       fields,
		"CollectionID": collectionID,
		"TotalCount":   totalCount,
		"TotalPages":   totalPages,
		"CurrentPage":  pageNum,
		"PageSize":     pageSizeNum,
	}, http.StatusOK)
}

func (ct *Controller) showEditContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	cID, ok := utils.StringToUint(c.Param("contentID"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	contentEntry, err := ct.services.Content.FindByID(cID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	collection, err := ct.services.Collection.FindByID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	contents, err := ct.services.Content.FindContentsWithDisplayContentValue()
	assets, _, err := ct.services.Asset.List(1, 100000)

	utils.RenderWithLayout(c, "content/create_or_edit.tmpl", gin.H{
		"FieldsHtml": RenderFieldsByContent(*contentEntry, *collection, contents, assets),
		"Collection": collection,
		"Content":    contentEntry,
	}, http.StatusOK)
}

func (ct *Controller) deleteContent(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("contentID"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	if err := ct.services.Content.DeleteByID(id); err != nil {
		ct.services.Webhook.Dispatch(string(model.EventContentDeleted), nil)
	}

	c.Redirect(http.StatusSeeOther, "/content/collections")
}
