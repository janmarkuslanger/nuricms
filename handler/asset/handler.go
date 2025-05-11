package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
)

type Handler struct {
	services *service.Set
}

func NewHandler(services *service.Set) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/assets", h.showAssets)
	r.GET("/assets/create", h.showCreateAsset)
	r.POST("/assets/create", h.createAsset)
	r.GET("/assets/edit/:id", h.showEditAsset)
	r.POST("/assets/edit/:id", h.editAsset)
	r.POST("/assets/delete/:id", h.deleteAsset)
}

func (h *Handler) showAssets(c *gin.Context) {
	assets, _ := h.services.Asset.GetAll()
	utils.RenderWithLayout(c, "/asset/index.tmpl", gin.H{
		"assets": assets,
	}, http.StatusOK)
}

func (h *Handler) showCreateAsset(c *gin.Context) {
	utils.RenderWithLayout(c, "/asset/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (h *Handler) deleteAsset(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	h.services.Asset.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/assets")
}

func (h *Handler) showEditAsset(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	asset, err := h.services.Asset.FindByID(id)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	utils.RenderWithLayout(c, "/asset/create_or_edit.tmpl", gin.H{
		"Asset": asset,
	}, http.StatusOK)
}

func (h *Handler) createAsset(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	name := c.PostForm("name")
	filePath, err := h.services.Asset.UploadFile(c, file)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	h.services.Asset.Create(&model.Asset{
		Path: filePath,
		Name: name,
	})

	c.Redirect(http.StatusSeeOther, "/assets")
}

func (h *Handler) editAsset(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	asset, err := h.services.Asset.FindByID(id)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	file, err := c.FormFile("file")

	if file != nil && err == nil {
		path, err := h.services.Asset.UploadFile(c, file)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/assets")
		}

		asset.Path = path
	}

	name := c.PostForm("name")
	asset.Name = name

	h.services.Asset.Save(asset)

	c.Redirect(http.StatusSeeOther, "/assets")
}
