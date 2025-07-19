package content

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
	s.Handle("GET /content/collections", ct.showCollections,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /content/collections/{id}/show", ct.listContent,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /content/collections/{id}/create", ct.showCreateContent,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /content/collections/{id}/create", ct.createContent,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /content/collections/{id}/edit/{contentID}", ct.showEditContent,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /content/collections/{id}/edit/{contentID}", ct.editContent,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("POST /content/collections/{id}/delete/{contentID}", ct.deleteContent,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)
}

func (ct Controller) showCollections(ctx server.Context) {
	handler.HandleList(ctx, ct.services.Collection, "content/collections.tmpl")
}

func (ct *Controller) showCreateContent(ctx server.Context) {
	collectionID, ok := utils.GetParamOrRedirect(ctx, "/content/collections", "id")
	if !ok {
		return
	}

	fields, errF := ct.services.Field.FindByCollectionID(collectionID)
	collection, errCol := ct.services.Collection.FindByID(collectionID)
	contents, errCon := ct.services.Content.FindContentsWithDisplayContentValue()
	assets, _, errA := ct.services.Asset.List(1, 100000)
	if errF != nil || errCol != nil || errCon != nil || errA != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/collections", http.StatusSeeOther)
		return
	}

	fieldsContent := make([]FieldContent, 0)
	for _, field := range fields {
		fieldsContent = append(fieldsContent, FieldContent{
			Field:   field,
			Content: contents,
			Assets:  assets,
		})
	}

	utils.RenderWithLayoutHTTP(ctx, "content/create_or_edit.tmpl", map[string]any{
		"FieldsHtml": RenderFields(fieldsContent),
		"Collection": collection,
	}, http.StatusOK)
}

func (ct *Controller) createContent(ctx server.Context) {
	collectionID, ok := utils.StringToUint(ctx.Request.PathValue("id"))
	if err := ctx.Request.ParseForm(); err != nil || !ok {
		http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
		return
	}

	if _, err := ct.services.Content.CreateWithValues(dto.ContentWithValues{
		CollectionID: collectionID,
		FormData:     ctx.Request.PostForm,
	}); err == nil {
		ct.services.Webhook.Dispatch(string(model.EventContentCreated), nil)
	}

	http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
}

func (ct *Controller) editContent(ctx server.Context) {
	colID, okCol := utils.StringToUint(ctx.Request.PathValue("id"))
	conID, okCon := utils.StringToUint(ctx.Request.PathValue("contentID"))
	err := ctx.Request.ParseForm()
	if !okCon || !okCol || err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
		return
	}

	if _, err := ct.services.Content.EditWithValues(dto.ContentWithValues{
		CollectionID: colID,
		ContentID:    conID,
		FormData:     ctx.Request.PostForm,
	}); err == nil {
		ct.services.Webhook.Dispatch(string(model.EventContentUpdated), nil)
	}

	http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
}

func (ct *Controller) listContent(ctx server.Context) {
	collectionID, ok := utils.GetParamOrRedirect(ctx, "/content/collections", "id")
	if !ok {
		return
	}

	page, pageSize := utils.ParsePagination(ctx.Request)

	contents, totalCount, err := ct.services.Content.FindDisplayValueByCollectionID(collectionID, page, pageSize)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
		return
	}

	groups := ContentsToContentGroup(contents)

	fields, err := ct.services.Field.FindDisplayFieldsByCollectionID(collectionID)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
		return
	}

	totalPages := utils.CalcTotalPages(totalCount, pageSize)

	utils.RenderWithLayoutHTTP(ctx, "content/content_list.tmpl", map[string]any{
		"Groups":       groups,
		"Fields":       fields,
		"CollectionID": collectionID,
		"TotalCount":   totalCount,
		"TotalPages":   totalPages,
		"CurrentPage":  page,
		"PageSize":     pageSize,
	}, http.StatusOK)
}

func (ct *Controller) showEditContent(ctx server.Context) {
	collectionID, okCol := utils.GetParamOrRedirect(ctx, "/content/collections", "id")
	cID, okCon := utils.GetParamOrRedirect(ctx, "/content/collections", "contentID")
	if !okCol || !okCon {
		return
	}

	contentEntry, err := ct.services.Content.FindByID(cID)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
		return
	}

	collection, err := ct.services.Collection.FindByID(collectionID)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
		return
	}

	contents, err := ct.services.Content.FindContentsWithDisplayContentValue()
	assets, _, err := ct.services.Asset.List(1, 100000)

	utils.RenderWithLayoutHTTP(ctx, "content/create_or_edit.tmpl", map[string]any{
		"FieldsHtml": RenderFieldsByContent(*contentEntry, DataContext{
			Collection: *collection,
			Contents:   contents,
			Assets:     assets,
		}),
		"Collection": collection,
		"Content":    contentEntry,
	}, http.StatusOK)
}

func (ct *Controller) deleteContent(ctx server.Context) {
	id, ok := utils.GetParamOrRedirect(ctx, "/content/collections", "contentID")
	if !ok {
		return
	}

	if err := ct.services.Content.DeleteByID(id); err == nil {
		ct.services.Webhook.Dispatch(string(model.EventContentDeleted), nil)
	}

	http.Redirect(ctx.Writer, ctx.Request, "/content/collections", http.StatusSeeOther)
}
