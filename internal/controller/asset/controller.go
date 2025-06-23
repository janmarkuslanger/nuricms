package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/assets", middleware.Userauth(ct.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showAssets)
	secure.GET("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showCreateAsset)
	secure.POST("/create", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.createAsset)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showEditAsset)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.editAsset)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), ct.deleteAsset)
}

func (h *Controller) showAssets(c *gin.Context) {
	page, pageSize := utils.ParsePagination(c)

	assets, totalCount, _ := h.services.Asset.List(page, pageSize)

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	utils.RenderWithLayout(c, "asset/index.tmpl", gin.H{
		"Assets":      assets,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": page,
		"PageSize":    pageSize,
	}, http.StatusOK)
}

func (ct *Controller) showCreateAsset(c *gin.Context) {
	utils.RenderWithLayout(c, "asset/create_or_edit.tmpl", gin.H{}, http.StatusOK)
}

func (ct *Controller) deleteAsset(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/assets")
		return
	}

	ct.services.Asset.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/assets")
}

func (ct *Controller) showEditAsset(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/assets")
		return
	}

	asset, err := ct.services.Asset.FindByID(id)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
		return
	}

	utils.RenderWithLayout(c, "asset/create_or_edit.tmpl", gin.H{
		"Asset": asset,
	}, http.StatusOK)
}

func (ct *Controller) createAsset(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
		return
	}

	name := c.PostForm("name")
	filePath, err := ct.services.Asset.UploadFile(c, file)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
		return
	}

	ct.services.Asset.Create(&model.Asset{
		Path: filePath,
		Name: name,
	})

	c.Redirect(http.StatusSeeOther, "/assets")
}

func (ct *Controller) editAsset(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/assets")
		return
	}

	asset, err := ct.services.Asset.FindByID(id)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/assets")
		return
	}

	file, err := c.FormFile("file")

	if file != nil && err == nil {
		path, err := ct.services.Asset.UploadFile(c, file)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/assets")
			return
		}

		asset.Path = path
	}

	name := c.PostForm("name")
	asset.Name = name

	ct.services.Asset.Save(asset)

	c.Redirect(http.StatusSeeOther, "/assets")
}
