package content

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/db"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
	"gorm.io/gorm"
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

	fields, err := ct.services.Field.GetByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	collection, err := ct.services.Collection.FindByID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	contents, err := ct.services.Content.GetContentsWithDisplayContentValue()
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
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	db.DB.Transaction(func(tx *gorm.DB) error {
		fields, err := ct.services.Field.GetByCollectionID(collectionID)
		if err != nil {
			return err
		}

		content := model.Content{CollectionID: collectionID}
		newContent, err := ct.services.Content.Create(&content)
		if err != nil {
			return err
		}

		for _, field := range fields {
			if field.IsList {
				vals := c.PostFormArray(field.Alias)
				for idx, val := range vals {
					cv := model.ContentValue{
						SortIndex: idx + 1,
						ContentID: newContent.ID,
						FieldID:   field.ID,
						Value:     val,
					}
					if err := ct.services.ContentValue.Create(&cv); err != nil {
						return err
					}
				}
			} else {
				val := c.PostForm(field.Alias)
				cv := model.ContentValue{
					SortIndex: 1,
					ContentID: newContent.ID,
					FieldID:   field.ID,
					Value:     val,
				}
				if err := ct.services.ContentValue.Create(&cv); err != nil {
					return err
				}
			}
		}

		return nil
	})

	ct.services.Webhook.Dispatch(string(model.EventContentCreated), nil)
	c.Redirect(http.StatusSeeOther, "/content/collections")
}

func (ct *Controller) editContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	id, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid Content-ID")
		return
	}
	contentID := uint(id)

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		existingContent, err := ct.services.Content.FindByID(contentID)
		if err != nil {
			return err
		}

		if existingContent.CollectionID != collectionID {
			return fmt.Errorf("Content %d doesnt relate to Collection %d", contentID, collectionID)
		}

		if err := tx.Where("content_id = ?", contentID).Delete(&model.ContentValue{}).Error; err != nil {
			return err
		}

		fields, err := ct.services.Field.GetByCollectionID(collectionID)
		if err != nil {
			return err
		}

		for _, field := range fields {
			if field.IsList {
				vals := c.PostFormArray(field.Alias)
				for idx, val := range vals {
					cv := model.ContentValue{
						SortIndex: idx + 1,
						ContentID: contentID,
						FieldID:   field.ID,
						Value:     val,
					}
					if err := tx.Create(&cv).Error; err != nil {
						return err
					}
				}
			} else {
				val := c.PostForm(field.Alias)
				cv := model.ContentValue{
					SortIndex: 1,
					ContentID: contentID,
					FieldID:   field.ID,
					Value:     val,
				}
				if err := tx.Create(&cv).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error while updating: %v", err)
		return
	}

	ct.services.Webhook.Dispatch(string(model.EventContentUpdated), nil)
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

	contents, totalCount, err := ct.services.Content.GetDisplayValueByCollectionID(collectionID, pageNum, pageSizeNum)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	groups := ContentsToContentGroup(contents)

	fields, err := ct.services.Field.GetDisplayFieldsByCollectionID(collectionID)
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

	contParam := c.Param("contentID")
	contID64, err := strconv.ParseUint(contParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}
	contentID := uint(contID64)

	contentEntry, err := ct.services.Content.FindByID(contentID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	collection, err := ct.services.Collection.FindByID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	contents, err := ct.services.Content.GetContentsWithDisplayContentValue()
	assets, _, err := ct.services.Asset.List(1, 100000)

	utils.RenderWithLayout(c, "content/create_or_edit.tmpl", gin.H{
		"FieldsHtml": RenderFieldsByContent(contentEntry, *collection, contents, assets),
		"Collection": collection,
		"Content":    contentEntry,
	}, http.StatusOK)
}

func (ct *Controller) deleteContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	id, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}
	contentID := uint(id)

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		content, err := ct.services.Content.FindByID(contentID)
		if err != nil {
			return err
		}
		if content.CollectionID != collectionID {
			return fmt.Errorf("Content %d doesnt relate to Collection %d", contentID, collectionID)
		}
		if err := tx.Where("content_id = ?", contentID).Delete(&model.ContentValue{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&content).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error while deleting: %v", err)
		return
	}

	ct.services.Webhook.Dispatch(string(model.EventContentDeleted), nil)
	c.Redirect(http.StatusSeeOther, "/content/collections")
}
