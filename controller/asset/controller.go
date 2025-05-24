package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (h *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/assets", middleware.Userauth(h.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showAssets)
	secure.GET("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showCreateAsset)
	secure.POST("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.createAsset)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showEditAsset)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.editAsset)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), h.deleteAsset)
}

func (h *Controller) showAssets(c *gin.Context) {
	assets, _ := h.services.Asset.GetAll()
	utils.RenderWithLayout(c, "/asset/index.tmpl", gin.H{
		"assets": assets,
	}, http.StatusOK)
}

func (h *Controller) showCreateAsset(c *gin.Context) {
	utils.RenderWithLayout(c, "/asset/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (h *Controller) deleteAsset(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/assets")
	}

	h.services.Asset.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/assets")
}

func (h *Controller) showEditAsset(c *gin.Context) {
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

func (h *Controller) createAsset(c *gin.Context) {
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

func (h *Controller) editAsset(c *gin.Context) {
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
