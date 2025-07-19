package apikey

import (
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/handler"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct Controller) RegisterRoutes(s *server.Server) {
	s.Handle("GET /apikeys",
		ct.showApikeys,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /apikeys/create",
		ct.showCreateApikey,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /apikeys/create",
		ct.createApikey,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("GET /apikeys/edit/{id}",
		ct.showEditApikey,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /apikeys/delete/{id}",
		ct.deleteApikey,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)
}

func (ct Controller) showApikeys(ctx server.Context) {
	handler.HandleList(ctx, ct.services.Apikey, "apikey/index.tmpl")
}

func (ct Controller) showCreateApikey(ctx server.Context) {
	handler.HandleShowCreate(ctx, handler.HandlerOptions{
		RenderOnSuccess: "apikey/create_or_edit.tmpl",
	})
}

func (ct Controller) createApikey(ctx server.Context) {
	handler.HandleCreate(ctx, ct.services.Apikey, dto.ApikeyData{
		Name: ctx.Request.PostFormValue("name"),
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/apikeys",
		RenderOnFail:      "apikey/create_or_edit.tmpl",
	})
}

func (ct Controller) showEditApikey(ctx server.Context) {
	handler.HandleShowEdit(ctx, ct.services.Apikey, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnFail:  "/apikeys",
		RenderOnSuccess: "apikey/create_or_edit.tmpl",
	})
}

func (ct Controller) deleteApikey(ctx server.Context) {
	handler.HandleDelete(ctx, ct.services.Apikey, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnSuccess: "/apikeys",
		RedirectOnFail:    "/apikeys",
	})
}
