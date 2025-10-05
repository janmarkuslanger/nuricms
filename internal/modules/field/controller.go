package field

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
	s.Handle("GET /fields",
		ct.listFields,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /fields/create",
		ct.showCreateField,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /fields/create",
		ct.createField,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("GET /fields/edit/{id}",
		ct.showEditField,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /fields/edit/{id}",
		ct.editField,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /fields/delete/{id}",
		ct.deleteField,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)
}

func (ct Controller) listFields(ctx server.Context) {
	handler.HandleList(ctx, ct.services.Field, "field/index.tmpl")
}

func loadFieldData(s service.CollectionService) (map[string]any, error) {
	data := make(map[string]any, 2)

	collections, _, err := s.List(1, 999999999999999999)
	if err != nil {
		return data, err
	}

	data["Collections"] = collections
	data["Types"] = []model.FieldType{
		model.FieldTypeText, model.FieldTypeNumber, model.FieldTypeBoolean,
		model.FieldTypeDate, model.FieldTypeAsset, model.FieldTypeCollection,
		model.FieldTypeTextarea, model.FieldTypeRichText, model.FieldTypeMultiSelect,
	}

	return data, nil
}

func (ct Controller) showCreateField(ctx server.Context) {
	handler.HandleShowCreate(ctx, handler.HandlerOptions{
		RenderOnSuccess: "field/create_or_edit.tmpl",
		TemplateData: func() (map[string]any, error) {
			return loadFieldData(ct.services.Collection)
		},
	})
}

func (ct Controller) createField(ctx server.Context) {
	handler.HandleCreate(ctx, ct.services.Field, dto.FieldData{
		Name:         ctx.Request.PostFormValue("name"),
		Alias:        ctx.Request.PostFormValue("alias"),
		CollectionID: ctx.Request.PostFormValue("collection_id"),
		FieldType:    ctx.Request.PostFormValue("field_type"),
		IsList:       ctx.Request.PostFormValue("is_list"),
		IsRequired:   ctx.Request.PostFormValue("is_required"),
		DisplayField: ctx.Request.PostFormValue("display_field"),
	}, handler.HandlerOptions{
		RenderOnFail:      "field/create_or_edit.tmpl",
		RedirectOnSuccess: "/fields/",
	})
}

func (ct Controller) showEditField(ctx server.Context) {
	handler.HandleShowEdit(ctx, ct.services.Field, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnFail:  "/fields/",
		RenderOnSuccess: "field/create_or_edit.tmpl",
		TemplateData: func() (map[string]any, error) {
			return loadFieldData(ct.services.Collection)
		},
	})
}

func (ct *Controller) editField(ctx server.Context) {
	handler.HandleEdit(ctx, ct.services.Field, ctx.Request.PathValue("id"), dto.FieldData{
		Name:         ctx.Request.PostFormValue("name"),
		Alias:        ctx.Request.PostFormValue("alias"),
		CollectionID: ctx.Request.PostFormValue("collection_id"),
		FieldType:    ctx.Request.PostFormValue("field_type"),
		IsList:       ctx.Request.PostFormValue("is_list"),
		IsRequired:   ctx.Request.PostFormValue("is_required"),
		DisplayField: ctx.Request.PostFormValue("display_field"),
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/fields/",
		RenderOnFail:      "field/create_or_edit.tmpl",
	})
}

func (ct *Controller) deleteField(ctx server.Context) {
	handler.HandleDelete(ctx, ct.services.Field, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnSuccess: "/fields/",
		RedirectOnFail:    "/fields/",
	})
}
