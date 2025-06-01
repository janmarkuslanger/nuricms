package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/middleware"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/user", middleware.Userauth(ct.services.User))

	r.GET("/login", ct.showLogin)
	r.POST("/login", ct.login)

	secure.GET("/", middleware.Roleauth(model.RoleEditor, model.RoleAdmin), ct.showUser)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), ct.showCreateUser)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), ct.createUser)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), ct.showEditUser)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleAdmin), ct.editUser)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), ct.deleteUser)
}

func (ct *Controller) showLogin(c *gin.Context) {
	utils.RenderWithLayout(c, "auth/login.tmpl", gin.H{}, http.StatusOK)
}

func (ct *Controller) login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	token, err := ct.services.User.LoginUser(email, password)
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
		gin.Mode() == gin.ReleaseMode,
		gin.Mode() == gin.ReleaseMode,
	)

	c.Redirect(http.StatusSeeOther, "/")
}

func (ct *Controller) showUser(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeNum = 10
	}

	users, totalCount, err := ct.services.User.List(pageNum, pageSizeNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users."})
		return
	}

	totalPages := (totalCount + int64(pageSizeNum) - 1) / int64(pageSizeNum)

	utils.RenderWithLayout(c, "user/index.tmpl", gin.H{
		"Roles":       model.GetUserRoles(),
		"User":        users,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": pageNum,
		"PageSize":    pageSizeNum,
	}, http.StatusOK)
}

func (ct *Controller) showCreateUser(c *gin.Context) {
	utils.RenderWithLayout(c, "user/create_or_edit.tmpl", gin.H{
		"Roles": model.GetUserRoles(),
	}, http.StatusOK)
}

func (ct *Controller) showEditUser(c *gin.Context) {
	userID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	user, err := ct.services.User.FindByID(userID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user")
	}

	utils.RenderWithLayout(c, "user/create_or_edit.tmpl", gin.H{
		"Roles": model.GetUserRoles(),
		"User":  user,
	}, http.StatusOK)
}

func (ct *Controller) editUser(c *gin.Context) {
	userID, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	user, err := ct.services.User.FindByID(userID)
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

	ct.services.User.Save(user)
	c.Redirect(http.StatusSeeOther, "/user")
}

func (ct *Controller) createUser(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("email")
	role := c.PostForm("role")

	_, err := ct.services.User.Create(email, password, model.Role(role))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/user")
}

func (ct *Controller) deleteUser(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/user")
	}

	ct.services.User.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/user")
}
