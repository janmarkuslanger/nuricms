package fieldoptions

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/dto"
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
	s.Handle("GET /field-options",
		ct.listFieldOptions,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /field-options/SelectOption",
		ct.indexFieldOptions,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /field-options/SelectOption/create",
		ct.showCreateFieldOption,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /field-options/SelectOption/create",
		ct.createFieldOption,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /field-options/SelectOption/edit/{id}",
		ct.showEditFieldOption,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /field-options/SelectOption/edit/{id}",
		ct.editFieldOption,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /field-options/SelectOption/delete/{id}",
		ct.deleteFieldOption,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)
}

func (ct Controller) listFieldOptions(ctx server.Context) {
	utils.RenderWithLayoutHTTP(ctx, "field_option/list.tmpl", map[string]any{
		"Types": []string{string(model.FieldOptionTypeSelectOption)},
	}, http.StatusOK)
}

func (ct Controller) indexFieldOptions(ctx server.Context) {
	handler.HandleList(ctx, ct.services.FieldOption, "field_option/index.tmpl")
}

func (ct Controller) showCreateFieldOption(ctx server.Context) {
	fields, _ := ct.services.Field.FindByFieldTypes([]model.FieldType{
		model.FieldTypeMultiSelect,
	})

	utils.RenderWithLayoutHTTP(ctx, "field_option/create_or_edit.tmpl", map[string]any{
		"Type":   string(model.FieldOptionTypeSelectOption),
		"Fields": fields,
	}, http.StatusOK)
}

func (ct Controller) createFieldOption(ctx server.Context) {
	dto := dto.FieldOption{
		Value:   ctx.Request.FormValue("value"),
		FieldID: ctx.Request.FormValue("field"),
	}

	ct.services.FieldOption.Create(dto)
	http.Redirect(ctx.Writer, ctx.Request, "/field-options/SelectOption", http.StatusSeeOther)
}

func (ct Controller) editFieldOption(ctx server.Context) {
	handler.HandleEdit(ctx, ct.services.FieldOption, ctx.Request.PathValue("id"), dto.FieldOption{
		Value: ctx.Request.PostFormValue("value"),
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/field-options/SelectOption",
		RedirectOnFail:    "/field-options/SelectOption/edit/" + ctx.Request.PathValue("id"),
	})
}

func (ct Controller) showEditFieldOption(ctx server.Context) {
	handler.HandleShowEdit(ctx, ct.services.FieldOption, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnFail:  "/field-options/SelectOption",
		RenderOnSuccess: "field_option/create_or_edit.tmpl",
	})
}

func (ct *Controller) deleteFieldOption(ctx server.Context) {
	handler.HandleDelete(ctx, ct.services.FieldOption, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnSuccess: "/field-options/SelectOption/",
		RedirectOnFail:    "/field-options/SelectOption/",
	})
}
