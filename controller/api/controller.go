package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
)

type Controller struct {
	services *service.Set
}

type CollectionResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type ContentResponse struct {
	Collection CollectionResponse    `json:"collection"`
	Items      []ContentItemResponse `json:"items"`
}

type ContentItemResponse struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Values    map[string]any `json:"values"`
}

type ContentValueResponse struct {
	ID    uint `json:"id"`
	Value any  `json:"value"`
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (h *Controller) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api",
		middleware.ApikeyAuth(h.services.Apikey),
	)
	api.GET("/collections/:alias/contents", h.listContents)
}
func (h *Controller) transformContentRecursive(ce *model.Content) ContentItemResponse {

	contentValues := make(map[string]any)

	for _, cv := range ce.ContentValues {
		alias := cv.Field.Alias
		var val any

		switch cv.Field.FieldType {
		case model.FieldTypeCollection:
			id, _ := strconv.ParseUint(cv.Value, 10, 32)
			cont, _ := h.services.Content.FindByID(uint(id))
			val = h.transformContentRecursive(&cont)

		case model.FieldTypeAsset:
			id, _ := strconv.ParseUint(cv.Value, 10, 32)
			asset, _ := h.services.Asset.FindByID(uint(id))
			val = asset

		default:
			val = cv.Value
		}

		if cv.Field.IsList {
			slice, ok := contentValues[alias].([]any)
			if !ok {
				slice = []any{}
			}
			slice = append(slice, val)
			contentValues[alias] = slice

		} else {
			contentValues[alias] = val
		}
	}

	return ContentItemResponse{
		ID:        ce.ID,
		CreatedAt: ce.CreatedAt,
		UpdatedAt: ce.UpdatedAt,
		Values:    contentValues,
	}
}

func (h *Controller) listContents(c *gin.Context) {
	alias := c.Param("alias")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	const perPage = 100
	offset := (page - 1) * perPage

	col, err := h.services.Collection.FindByAlias(alias)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	content, err := h.services.Content.ListByCollectionAlias(alias, offset, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var contentItems []ContentItemResponse

	for _, ce := range content {
		contentItems = append(contentItems, h.transformContentRecursive(&ce))
	}

	out := ContentResponse{
		Collection: CollectionResponse{
			ID:    col.ID,
			Name:  col.Name,
			Alias: col.Alias,
		},
		Items: contentItems,
	}

	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"perPage": perPage,
		"data":    out,
	})
}

// func (h *Controller) transformContentRecursive(ct *model.Content) ContentResponseMap {
// 	values := make(map[string]ContentField)

// 	for _, cv := range ct.ContentValues {
// 		var val any = cv.Value

// 		switch cv.Field.FieldType {
// 		case "Asset":
// 			if id, err := strconv.ParseUint(cv.Value, 10, 32); err == nil {
// 				if asset, err2 := h.services.Asset.FindByID(uint(id)); err2 == nil {
// 					val = asset
// 				}
// 			}
// 		case "Collection":
// 			if id, err := strconv.ParseUint(cv.Value, 10, 32); err == nil {
// 				if nested, err2 := h.services.Content.FindByID(uint(id)); err2 == nil {
// 					val = h.transformContentRecursive(&nested)
// 				}
// 			}
// 		}

// 		cf, ok := values[cv.Field.Alias]
// 		item := ContentFieldValue{
// 			Collection: ct.Collection,
// 			Type:       string(cv.Field.FieldType),
// 			Value:      val,
// 		}

// 		if !ok {
// 			var items []ContentFieldValue
// 			items = append(items, item)
// 			values[cv.Field.Alias] = ContentField{
// 				ID:    cv.ID,
// 				Items: items,
// 				Field: Field{
// 					Name:   cv.Field.Name,
// 					Alias:  cv.Field.Alias,
// 					IsList: cv.Field.IsList,
// 				},
// 			}
// 		} else {
// 			cf.Items = append(cf.Items, item)
// 			values[cv.Field.Alias] = cf
// 		}

// 	}

// 	return ContentResponseMap{
// 		ID:     ct.ID,
// 		Values: values,
// 	}
// }
