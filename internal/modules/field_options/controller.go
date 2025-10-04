package fieldoptions

import (
	"net/http"

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
}

func (ct Controller) listFieldOptions(ctx server.Context) {
	utils.RenderWithLayoutHTTP(ctx, "field_option/list.tmpl", map[string]any{
		"Types": []string{string(model.FieldOptionTypeSelectOption)},
	}, http.StatusOK)
}

func (ct Controller) indexFieldOptions(ctx server.Context) {
	utils.RenderWithLayoutHTTP(ctx, "field_option/index.tmpl", map[string]any{
		"Type": string(model.FieldOptionTypeSelectOption),
	}, http.StatusOK)
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
	v := ctx.Request.FormValue("value")
	f := ctx.Request.FormValue("field")

	fo := model.FieldOption{
		OptionType: model.FieldOptionTypeSelectOption,
		Value:      v,
		FieldID:    uint(f),
	}

	utils.RenderWithLayoutHTTP(ctx, "field_option/create_or_edit.tmpl", map[string]any{}, http.StatusOK)
}
