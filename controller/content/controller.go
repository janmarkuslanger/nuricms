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

func (h *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/content/collections", middleware.Userauth(h.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showCollections)
	secure.GET("/:id/show", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.listContent)
	secure.GET("/:id/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showCreateContent)
	secure.POST("/:id/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.createContent)
	secure.GET("/:id/edit/:contentID", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showEditContent)
	secure.POST("/:id/edit/:contentID", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.editContent)
	secure.POST("/:id/delete/:contentID", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.deleteContent)
}

func (h *Controller) showCollections(c *gin.Context) {
	collections, err := h.services.Collection.GetAll()
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}
	utils.RenderWithLayout(c, "content/collections.tmpl", gin.H{
		"collections": collections,
	}, http.StatusOK)
}

func (h *Controller) showCreateContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	fields, err := h.services.Field.GetByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	collection, err := h.services.Collection.GetByID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	contents, err := h.services.Content.GetContentsWithDisplayContentValue()
	assets, err := h.services.Asset.GetAll()

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

func (h *Controller) createContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	db.DB.Transaction(func(tx *gorm.DB) error {
		fields, err := h.services.Field.GetByCollectionID(collectionID)
		if err != nil {
			return err
		}

		content := model.Content{CollectionID: collectionID}
		newContent, err := h.services.Content.Create(&content)
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
					if err := h.services.ContentValue.Create(&cv); err != nil {
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
				if err := h.services.ContentValue.Create(&cv); err != nil {
					return err
				}
			}
		}

		return nil
	})

	c.Redirect(http.StatusSeeOther, "/content/collections")
}

func (h *Controller) editContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	id, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.String(http.StatusBadRequest, "Ungültige Content-ID")
		return
	}
	contentID := uint(id)

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		existingContent, err := h.services.Content.FindByID(contentID)
		if err != nil {
			return err
		}

		if existingContent.CollectionID != collectionID {
			return fmt.Errorf("Content %d doesnt relate to Collection %d", contentID, collectionID)
		}

		if err := tx.Where("content_id = ?", contentID).Delete(&model.ContentValue{}).Error; err != nil {
			return err
		}

		fields, err := h.services.Field.GetByCollectionID(collectionID)
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

	c.Redirect(http.StatusSeeOther, "/content/collections")
}

func (h *Controller) listContent(c *gin.Context) {
	collectionID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	contents, err := h.services.Content.GetDisplayValueByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	groups := ContentsToContentGroup(contents)

	fields, err := h.services.Field.GetDisplayFieldsByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	utils.RenderWithLayout(c, "content/content_list.tmpl", gin.H{
		"Groups":       groups,
		"Fields":       fields,
		"CollectionID": collectionID,
	}, http.StatusOK)
}

func (h *Controller) showEditContent(c *gin.Context) {
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

	contentEntry, err := h.services.Content.FindByID(contentID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	collection, err := h.services.Collection.GetByID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	contents, err := h.services.Content.GetContentsWithDisplayContentValue()
	assets, err := h.services.Asset.GetAll()

	utils.RenderWithLayout(c, "content/create_or_edit.tmpl", gin.H{
		"FieldsHtml": RenderFieldsByContent(contentEntry, *collection, contents, assets),
		"Collection": collection,
		"Content":    contentEntry,
	}, http.StatusOK)
}

func (h *Controller) deleteContent(c *gin.Context) {
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
		content, err := h.services.Content.FindByID(contentID)
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

	c.Redirect(http.StatusSeeOther, "/content/collections")
}
