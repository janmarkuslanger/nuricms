package collection

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
	s.Handle("GET /collections",
		ct.showCollections,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /collections/create",
		ct.showCreateCollection,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /collections/create",
		ct.createCollection,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("GET /collections/edit/{id}",
		ct.showEditCollection,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /collections/edit/{id}",
		ct.editCollection,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /collections/delete/{id}",
		ct.deleteCollection,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)
}

func (ct Controller) showCollections(ctx server.Context) {
	handler.HandleList(ctx, ct.services.Collection, "collection/index.tmpl")
}

func (ct Controller) showCreateCollection(ctx server.Context) {
	handler.HandleShowCreate(ctx, handler.HandlerOptions{
		RenderOnSuccess: "collection/create_or_edit.tmpl",
	})
}

func (ct Controller) createCollection(ctx server.Context) {
	handler.HandleCreate(ctx, ct.services.Collection, dto.CollectionData{
		Name:        ctx.Request.PostFormValue("name"),
		Alias:       ctx.Request.PostFormValue("alias"),
		Description: ctx.Request.PostFormValue("description"),
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/collections",
		RenderOnFail:      "collection/create_or_edit.tmpl",
	})
}

func (ct Controller) showEditCollection(ctx server.Context) {
	handler.HandleShowEdit(ctx, ct.services.Collection, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnFail:  "/collections/",
		RenderOnSuccess: "collection/create_or_edit.tmpl",
	})
}

func (ct Controller) editCollection(ctx server.Context) {
	handler.HandleEdit(ctx, ct.services.Collection, ctx.Request.PathValue("id"), dto.CollectionData{
		Name:        ctx.Request.PostFormValue("name"),
		Alias:       ctx.Request.PostFormValue("alias"),
		Description: ctx.Request.PostFormValue("description"),
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/collections/",
		RenderOnFail:      "collection/create_or_edit.tmpl",
	})
}

func (ct Controller) deleteCollection(ctx server.Context) {
	handler.HandleDelete(ctx, ct.services.Collection, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnSuccess: "/collections/",
		RedirectOnFail:    "/collections/",
	})
}
