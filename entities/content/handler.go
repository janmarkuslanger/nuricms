package content

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/core/db"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
	"github.com/janmarkuslanger/nuricms/utils"
	"gorm.io/gorm"
)

type Handler struct {
	repos *repository.Set
}

func NewHandler(repos *repository.Set) *Handler {
	return &Handler{repos: repos}
}

func (h Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/content/collections", h.showCollections)
	r.GET("/content/collections/:id/show", h.listContent)

	r.GET("/content/collections/:id/create", h.showCreateContent)
	r.POST("/content/collections/:id/create", h.createContent)

	r.GET("/content/collections/:id/edit/:contentID", h.showEditContent)
}

func (h Handler) showCollections(c *gin.Context) {
	collections, err := h.repos.Collection.GetAll()

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	utils.RenderWithLayout(c, "content/collections.html", gin.H{
		"collections": collections,
	}, http.StatusOK)
}

func (h Handler) showCreateContent(c *gin.Context) {
	collectionIDParam := c.Param("id")
	collectionID64, err := strconv.ParseUint(collectionIDParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	collectionID := uint(collectionID64)

	fields, err := h.repos.Field.FindByCollectionID(collectionID)

	collection, err := h.repos.Collection.FindByID(collectionID)

	utils.RenderWithLayout(c, "content/create_or_edit.html", gin.H{
		"FieldsHtml": RenderFields(fields),
		"Collection": collection,
	}, http.StatusOK)
}

func (h Handler) createContent(c *gin.Context) {
	collectionIDParam := c.Param("id")
	collectionID64, err := strconv.ParseUint(collectionIDParam, 10, 0)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	collectionID := uint(collectionID64)

	db.DB.Transaction(func(tx *gorm.DB) error {
		fields, err := h.repos.Field.FindByCollectionID(collectionID)

		if err != nil {
			return err
		}

		content := model.Content{
			CollectionID: collectionID,
		}

		newContent, err := h.repos.Content.Create(&content)

		if err != nil {
			return err
		}

		for _, field := range fields {

			if field.IsList {

				fieldValues := c.PostFormArray(field.Alias)

				for index, fieldValue := range fieldValues {
					contentValue := model.ContentValue{
						SortIndex: int(index + 1),
						ContentID: newContent.ID,
						FieldID:   field.ID,
						Value:     fieldValue,
					}
					_, err := h.repos.ContentValue.Create(&contentValue)

					if err != nil {
						return err
					}

				}
			} else {
				fieldValue := c.PostForm(field.Alias)
				contentValue := model.ContentValue{
					SortIndex: 1,
					ContentID: newContent.ID,
					FieldID:   field.ID,
					Value:     fieldValue,
				}
				_, err := h.repos.ContentValue.Create(&contentValue)

				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	c.Redirect(http.StatusSeeOther, "/content/collections")
}

func (h Handler) listContent(c *gin.Context) {
	collectionIDParam := c.Param("id")
	collectionID64, err := strconv.ParseUint(collectionIDParam, 10, 0)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	collectionID := uint(collectionID64)

	contents, err := h.repos.Content.FindDisplayValueByCollectionID(collectionID)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/content/collections")
		return
	}

	groups := GroupContentValuesByField(contents)

	fields, err := h.repos.Field.FindDisplayFieldsByCollectionID(collectionID)

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
	collectionIDParam := c.Param("id")
	collectionID64, err := strconv.ParseUint(collectionIDParam, 10, 32)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	collectionID := uint(collectionID64)

	fields, err := h.repos.Field.FindByCollectionID(collectionID)

	collection, err := h.repos.Collection.FindByID(collectionID)

	utils.RenderWithLayout(c, "content/create_or_edit.html", gin.H{
		"FieldsHtml": RenderFields(fields),
		"Collection": collection,
	}, http.StatusOK)
}
