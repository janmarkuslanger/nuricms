package content

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type Handler struct {
	repos *repository.Set
}

func NewHandler(repos *repository.Set) *Handler {
	return &Handler{repos: repos}
}

func (h Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/content/collections", h.showCollections)
	r.GET("/content/collections/create/:id", h.showCreateCollection)
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

func (h Handler) showCreateCollection(c *gin.Context) {
	collectionIDParam := c.Param("id")
	collectionID, err := strconv.ParseUint(collectionIDParam, 10, 32)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/collections")
		return
	}

	fields, err := h.repos.Field.FindByCollectionID(collectionID)

	utils.RenderWithLayout(c, "content/create_or_edit.html", gin.H{
		"fieldsHtml": RenderFields(fields),
	}, http.StatusOK)
}
