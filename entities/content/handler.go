package content

import (
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
}

func (h *Handler) showCollections(c *gin.Context) {
	collections, err := h.services.Collection.GetAll()
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}
	utils.RenderWithLayout(c, "content/collections.html", gin.H{
		"collections": collections,
	}, http.StatusOK)
}

func (h *Handler) showCreateContent(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}
	collectionID := uint(id64)

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

	utils.RenderWithLayout(c, "content/create_or_edit.html", gin.H{
		"FieldsHtml": RenderFields(fields),
		"Collection": collection,
	}, http.StatusOK)
}

func (h *Handler) createContent(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}
	collectionID := uint(id64)

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

func (h *Handler) listContent(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}
	collectionID := uint(id64)

	contents, err := h.services.Content.GetDisplayValueByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	groups := CreateContentGroupsByField(contents)

	fields, err := h.services.Field.GetDisplayFieldsByCollectionID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	utils.RenderWithLayout(c, "content/content_list.html", gin.H{
		"Groups":       groups,
		"Fields":       fields,
		"CollectionID": collectionID,
	}, http.StatusOK)
}

func (h *Handler) showEditContent(c *gin.Context) {
	collParam := c.Param("id")
	collID64, err := strconv.ParseUint(collParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}
	collectionID := uint(collID64)

	contParam := c.Param("contentID")
	contID64, err := strconv.ParseUint(contParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}
	contentID := uint(contID64)

	contentObj, err := h.services.Content.GetByID(contentID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	collection, err := h.services.Collection.GetByID(collectionID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	utils.RenderWithLayout(c, "content/create_or_edit.html", gin.H{
		"FieldsHtml": RenderFieldsByContent(contentObj),
		"Collection": collection,
	}, http.StatusOK)
}
