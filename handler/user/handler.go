package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
)

type Handler struct {
	services *service.Set
}

func NewHandler(services *service.Set) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/user", middleware.Userauth(h.services.User))

	r.GET("/login", h.showLogin)
	r.POST("/login", h.login)

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), h.showUser)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), h.showCreateUser)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), h.createUser)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), h.showEditUser)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleAdmin), h.editUser)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), h.deleteUser)
}

func (h *Handler) showLogin(c *gin.Context) {
	utils.RenderWithLayout(c, "auth/login.tmpl", gin.H{}, http.StatusOK)
}

func (h *Handler) login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	token, err := h.services.User.LoginUser(email, password)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		3600*24,
		"/",
		"",
		true,
		true,
	)

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) showUser(c *gin.Context) {
	user, _ := h.services.User.List()

	utils.RenderWithLayout(c, "user/index.tmpl", gin.H{
		"Roles": model.GetUserRoles(),
		"User":  user,
	}, http.StatusOK)
}

func (h *Handler) showCreateUser(c *gin.Context) {
	utils.RenderWithLayout(c, "user/create_or_edit.tmpl", gin.H{
		"Roles": model.GetUserRoles(),
	}, http.StatusOK)
}

func (h *Handler) showEditUser(c *gin.Context) {
	userID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	user, err := h.services.User.FindByID(userID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user")
	}

	utils.RenderWithLayout(c, "user/create_or_edit.tmpl", gin.H{
		"Roles": model.GetUserRoles(),
		"User":  user,
	}, http.StatusOK)
}

func (h *Handler) editUser(c *gin.Context) {
	userID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	user, err := h.services.User.FindByID(userID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user")
	}

	email := c.PostForm("email")
	password := c.PostForm("password")
	role := c.PostForm("role")

	user.Email = email
	user.Role = model.Role(role)

	if password != "" {
		user.Password = password
	}

	h.services.User.Save(user)
	c.Redirect(http.StatusSeeOther, "/user")
}

func (h *Handler) createUser(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("email")
	role := c.PostForm("role")

	_, err := h.services.User.Create(email, password, model.Role(role))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/user")
}

func (h *Handler) deleteUser(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/user")
	}

	h.services.User.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/user")
}
