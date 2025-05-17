// handlers/api_content.go
package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/service"
)

type ContentValueResponse struct {
	Field string      `json:"field"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type ContentResponse struct {
	ID     uint                   `json:"id"`
	Values []ContentValueResponse `json:"values"`
}

type Handler struct {
	services *service.Set
}

func NewHandler(services *service.Set) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/collections/:alias/contents", h.listContents)
}

// GET /api/collections/:alias/contents?page=1
func (h *Handler) listContents(c *gin.Context) {
	alias := c.Param("alias")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	const perPage = 100
	offset := (page - 1) * perPage

	contents, err := h.services.Content.
		ListByCollectionAlias(alias, offset, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var out []ContentResponse
	for _, ct := range contents {
		var cvres []ContentValueResponse
		for _, cv := range ct.ContentValues {
			var val interface{} = cv.Value

			switch cv.Field.FieldType {
			case "Asset":
				if id, err := strconv.ParseUint(cv.Value, 10, 32); err == nil {
					if asset, err2 := h.services.Asset.FindByID(uint(id)); err2 == nil {
						val = asset
					}
				}
			case "Collection":
				if id, err := strconv.ParseUint(cv.Value, 10, 32); err == nil {
					if nested, err2 := h.services.Content.FindByID(uint(id)); err2 == nil {
						val = nested
					}
				}
			}

			cvres = append(cvres, ContentValueResponse{
				Field: cv.Field.Name,
				Type:  string(cv.Field.FieldType),
				Value: val,
			})
		}
		out = append(out, ContentResponse{
			ID:     ct.ID,
			Values: cvres,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"perPage": perPage,
		"data":    out,
	})
}
