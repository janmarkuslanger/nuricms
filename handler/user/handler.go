package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	r.GET("/user", h.showUser)
	r.GET("/user/create", h.showCreateUser)
	r.POST("/user/create", h.createUser)
	r.GET("/user/edit/:id", h.showEditUser)
	r.POST("/user/edit/:id", h.editUser)
	r.POST("/user/delete/:id", h.deleteUser)
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
	password := c.PostForm("email")
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

	_, err := h.services.User.CreateUser(email, password, model.Role(role))
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
