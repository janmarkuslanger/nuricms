package asset

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/handler"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct Controller) RegisterRoutes(s *server.Server) {
	s.Handle("GET /assets",
		ct.showAssets,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /assets/create",
		ct.showCreateAsset,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /assets/create",
		ct.createAsset,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /assets/edit/{id}",
		ct.showEditAsset,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /assets/edit/{id}",
		ct.editAsset,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /assets/delete/{id}",
		ct.deleteAsset,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)
}

func (ct Controller) showAssets(ctx server.Context) {
	handler.HandleList(ctx, ct.services.Asset, "asset/index.tmpl")
}

func (ct Controller) showCreateAsset(ctx server.Context) {
	utils.RenderWithLayoutHTTP(ctx, "asset/create_or_edit.tmpl", map[string]any{}, http.StatusOK)
}

func (ct Controller) deleteAsset(ctx server.Context) {
	id, ok := utils.GetParamOrRedirect(ctx, "/assets", "id")
	if !ok {
		return
	}

	ct.services.Asset.DeleteByID(id)
	http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
}

func (ct Controller) showEditAsset(ctx server.Context) {
	id, ok := utils.GetParamOrRedirect(ctx, "/assets", "id")
	if !ok {
		return
	}

	asset, err := ct.services.Asset.FindByID(id)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
		return
	}

	utils.RenderWithLayoutHTTP(ctx, "asset/create_or_edit.tmpl", map[string]any{
		"Asset": asset,
	}, http.StatusOK)
}

func (ct Controller) createAsset(ctx server.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
		return
	}
	defer file.Close()

	name := ctx.Request.FormValue("name")
	filePath, err := ct.services.Asset.UploadFile(ctx, header)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
		return
	}

	ct.services.Asset.Create(&model.Asset{
		Path: filePath,
		Name: name,
	})

	http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
}

func (ct Controller) editAsset(ctx server.Context) {
	id, ok := utils.GetParamOrRedirect(ctx, "/assets", "id")
	if !ok {
		return
	}

	asset, err := ct.services.Asset.FindByID(id)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
		return
	}

	err = ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err == nil && file != nil {
		defer file.Close()
		path, err := ct.services.Asset.UploadFile(ctx, header)
		if err != nil {
			http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
			return
		}
		asset.Path = path
	}

	name := ctx.Request.FormValue("name")
	asset.Name = name

	ct.services.Asset.Save(asset)

	http.Redirect(ctx.Writer, ctx.Request, "/assets", http.StatusSeeOther)
}
