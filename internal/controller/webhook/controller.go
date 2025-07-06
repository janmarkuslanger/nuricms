package webhook

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

func (ct *Controller) RegisterRoutes(s *server.Server) {
	s.Handle("GET /webhooks",
		ct.showWebhooks,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /webhooks/create",
		ct.showCreateWebhook,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /webhooks/create",
		ct.createWebhook,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("GET /webhooks/edit/{id}",
		ct.showEditWebhook,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /webhooks/edit/{id}",
		ct.editWebhook,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /webhooks/delete/{id}",
		ct.deleteWebhook,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)
}

func (ct Controller) showWebhooks(ctx server.Context) {
	handler.HandleList(ctx, ct.services.Webhook, "webhook/index.tmpl")
}

func (ct Controller) showCreateWebhook(ctx server.Context) {
	handler.HandleShowCreate(ctx, handler.HandlerOptions{
		RenderOnSuccess: "webhook/create_or_edit.tmpl",
		TemplateData: func() (map[string]any, error) {
			data := make(map[string]any, 1)
			data["RequestTypes"] = model.GetRequestTypes()
			data["EventTypes"] = model.GetWebhookEvents()
			return data, nil
		},
	})
}

func (ct Controller) createWebhook(ctx server.Context) {
	events := make(map[string]bool)
	for _, event := range model.GetWebhookEvents() {
		events[string(event)] = ctx.Request.PostFormValue(string(event)) == "on"
	}

	handler.HandleCreate(ctx, ct.services.Webhook, dto.WebhookData{
		Name:        ctx.Request.PostFormValue("name"),
		Url:         ctx.Request.PostFormValue("url"),
		RequestType: ctx.Request.PostFormValue("request_type"),
		Events:      events,
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/webhooks",
		RenderOnFail:      "webhook/create_or_edit.tmpl",
	})
}

func (ct Controller) showEditWebhook(ctx server.Context) {
	handler.HandleShowEdit(ctx, ct.services.Webhook, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnFail:  "/webhooks",
		RenderOnSuccess: "webhook/create_or_edit.tmpl",
		TemplateData: func() (map[string]any, error) {
			data := make(map[string]any, 1)
			data["RequestTypes"] = model.GetRequestTypes()
			data["EventTypes"] = model.GetWebhookEvents()
			return data, nil
		},
	})
}

func (ct Controller) editWebhook(ctx server.Context) {
	events := make(map[string]bool)
	for _, event := range model.GetWebhookEvents() {
		events[string(event)] = ctx.Request.PostFormValue(string(event)) == "on"
	}

	handler.HandleEdit(ctx, ct.services.Webhook, ctx.Request.PathValue("id"), dto.WebhookData{
		Name:        ctx.Request.PostFormValue("name"),
		Url:         ctx.Request.PostFormValue("url"),
		RequestType: ctx.Request.PostFormValue("request_type"),
		Events:      events,
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/webhooks",
		RenderOnFail:      "webhook/create_or_edit.tmpl",
	})
}

func (ct Controller) deleteWebhook(ctx server.Context) {
	handler.HandleDelete(ctx, ct.services.Webhook, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnSuccess: "/webhooks",
		RedirectOnFail:    "/webhooks",
	})
}
