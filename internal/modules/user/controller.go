package user

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

func (ct *Controller) RegisterRoutes(s *server.Server) {
	s.Handle("GET /login", ct.showLogin)
	s.Handle("POST /login", ct.login)

	s.Handle("GET /user",
		ct.showUser,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin, model.RoleEditor),
	)

	s.Handle("GET /user/create",
		ct.showCreateUser,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /user/create",
		ct.createUser,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("GET /user/edit/{id}",
		ct.showEditUser,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /user/edit/{id}",
		ct.editUser,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)

	s.Handle("POST /user/delete/{id}",
		ct.deleteUser,
		middleware.Userauth(ct.services.User),
		middleware.Roleauth(model.RoleAdmin),
	)
}

func (ct *Controller) showLogin(ctx server.Context) {
	utils.RenderWithLayoutHTTP(ctx, "auth/login.tmpl", map[string]any{}, http.StatusOK)
}

func (ct *Controller) login(ctx server.Context) {
	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")

	token, err := ct.services.User.LoginUser(email, password)
	if err != nil {
		http.Redirect(ctx.Writer, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600 * 24,
	})

	http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
}

func (ct Controller) showUser(ctx server.Context) {
	handler.HandleList(ctx, ct.services.User, "user/index.tmpl")
}

func (ct Controller) showCreateUser(ctx server.Context) {
	handler.HandleShowCreate(ctx, handler.HandlerOptions{
		RenderOnSuccess: "user/create_or_edit.tmpl",
		TemplateData: func() (map[string]any, error) {
			data := make(map[string]any, 1)
			data["Roles"] = model.GetUserRoles()
			return data, nil
		},
	})
}

func (ct Controller) createUser(ctx server.Context) {
	handler.HandleCreate(ctx, ct.services.User, dto.UserData{
		Email:    ctx.Request.PostFormValue("email"),
		Password: ctx.Request.PostFormValue("password"),
		Role:     ctx.Request.PostFormValue("role"),
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/user",
		RenderOnFail:      "user/create_or_edit.tmpl",
	})
}

func (ct Controller) showEditUser(ctx server.Context) {
	handler.HandleShowEdit(ctx, ct.services.User, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnFail:  "/user",
		RenderOnSuccess: "user/create_or_edit.tmpl",
		TemplateData: func() (map[string]any, error) {
			data := make(map[string]any, 1)
			data["Roles"] = model.GetUserRoles()
			return data, nil
		},
	})
}

func (ct Controller) editUser(ctx server.Context) {
	handler.HandleEdit(ctx, ct.services.User, ctx.Request.PathValue("id"), dto.UserData{
		Email:    ctx.Request.PostFormValue("email"),
		Password: ctx.Request.PostFormValue("password"),
		Role:     ctx.Request.PostFormValue("role"),
	}, handler.HandlerOptions{
		RedirectOnSuccess: "/user",
		RenderOnFail:      "user/create_or_edit.tmpl",
	})
}

func (ct Controller) deleteUser(ctx server.Context) {
	handler.HandleDelete(ctx, ct.services.User, ctx.Request.PathValue("id"), handler.HandlerOptions{
		RedirectOnSuccess: "/user",
		RedirectOnFail:    "/user",
	})
}
