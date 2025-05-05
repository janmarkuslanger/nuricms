package content

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/db"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
	"gorm.io/gorm"
)

type Handler struct {
	services *service.Set
}

func NewHandler(services *service.Set) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/content/collections", h.showCollections)
	r.GET("/content/collections/:id/show", h.listContent)
	r.GET("/content/collections/:id/create", h.showCreateContent)
	r.POST("/content/collections/:id/create", h.createContent)
	r.GET("/content/collections/:id/edit/:contentID", h.showEditContent)
	r.POST("/content/collections/:id/edit/:contentID", h.editContent)
	r.POST("/content/collections/:id/delete/:contentID", h.deleteContent)
}

func (h *Handler) showCollections(c *gin.Context) {
	collections, err := h.services.Collection.GetAll()
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}
	utils.RenderWithLayout(c, "content/collections.tpl", gin.H{
		"collections": collections,
	}, http.StatusOK)
}

func (h *Handler) showCreateContent(c *gin.Context) {
	collectionID, ok := ParseCollectionID(c)
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

	fieldsContent := make([]FieldContent, 0)
	for _, field := range fields {
		fieldsContent = append(fieldsContent, FieldContent{
			Field:   field,
			Content: contents,
		})
	}

	utils.RenderWithLayout(c, "content/create_or_edit.tpl", gin.H{
		"FieldsHtml": RenderFields(fieldsContent),
		"Collection": collection,
	}, http.StatusOK)
}

func (h *Handler) createContent(c *gin.Context) {
	collectionID, ok := ParseCollectionID(c)
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

func (h *Handler) editContent(c *gin.Context) {
	collectionID, ok := ParseCollectionID(c)
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
		existingContent, err := h.services.Content.GetByID(contentID)
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

func (h *Handler) listContent(c *gin.Context) {
	collectionID, ok := ParseCollectionID(c)
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

	utils.RenderWithLayout(c, "content/content_list.tpl", gin.H{
		"Groups":       groups,
		"Fields":       fields,
		"CollectionID": collectionID,
	}, http.StatusOK)
}

func (h *Handler) showEditContent(c *gin.Context) {
	collectionID, ok := ParseCollectionID(c)
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

	contentEntry, err := h.services.Content.GetByID(contentID)
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

	utils.RenderWithLayout(c, "content/create_or_edit.tpl", gin.H{
		"FieldsHtml": RenderFieldsByContent(contentEntry, *collection, contents),
		"Collection": collection,
		"Content":    contentEntry,
	}, http.StatusOK)
}

func (h *Handler) deleteContent(c *gin.Context) {
	collectionID, ok := ParseCollectionID(c)
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
		content, err := h.services.Content.GetByID(contentID)
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
