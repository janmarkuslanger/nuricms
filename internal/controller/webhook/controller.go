package webhook

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func (ct *Controller) RegisterRoutes(r *gin.Engine) {
	secure := r.Group("/webhooks", middleware.Userauth(ct.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleAdmin), ct.showWebhooks)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), ct.showCreateWebhook)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), ct.createWebhook)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), ct.showEditWebhook)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleAdmin), ct.editWebhook)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), ct.deleteWebhook)
}

func (ct *Controller) showWebhooks(c *gin.Context) {
	page, pageSize := utils.ParsePagination(c)

	webhooks, totalCount, err := ct.services.Webhook.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve webhooks."})
		return
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	utils.RenderWithLayout(c, "webhook/index.tmpl", gin.H{
		"Webhooks":    webhooks,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": page,
		"PageSize":    pageSize,
	}, http.StatusOK)
}

func (ct *Controller) showCreateWebhook(c *gin.Context) {
	utils.RenderWithLayout(c, "webhook/create_or_edit.tmpl", gin.H{
		"RequestTypes": model.GetRequestTypes(),
		"EventTypes":   model.GetWebhookEvents(),
	}, http.StatusOK)
}

func (ct *Controller) createWebhook(c *gin.Context) {
	name := c.PostForm("name")
	url := c.PostForm("url")
	requestType := model.RequestType(c.PostForm("request_type"))

	events := map[model.EventType]bool{}

	for _, event := range model.GetWebhookEvents() {
		isActive := c.PostForm(string(event)) == "on"
		events[event] = isActive
	}

	_, err := ct.services.Webhook.Create(name, url, requestType, events)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/webhooks")
		return
	}

	c.Redirect(http.StatusSeeOther, "/webhooks")
}

func (ct *Controller) showEditWebhook(c *gin.Context) {
	id, ok := utils.GetParamOrRedirect(c, "/webhooks", "id")
	if !ok {
		return
	}

	webhook, err := ct.services.Webhook.FindByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	utils.RenderWithLayout(c, "webhook/create_or_edit.tmpl", gin.H{
		"RequestTypes":    model.GetRequestTypes(),
		"EventTypes":      model.GetWebhookEvents(),
		"Webhook":         webhook,
		"EventTypeValues": strings.Split(webhook.Events, ","),
	}, http.StatusOK)
}

func (ct *Controller) editWebhook(c *gin.Context) {
	id, ok := utils.GetParamOrRedirect(c, "/webhooks", "id")
	if !ok {
		return
	}

	webhook, err := ct.services.Webhook.FindByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	name := c.PostForm("name")
	url := c.PostForm("url")
	requestType := model.RequestType(c.PostForm("request_type"))

	var eventString strings.Builder
	for _, event := range model.GetWebhookEvents() {
		s := string(event)
		if c.PostForm(s) == "on" {
			eventString.WriteString(s)
			eventString.WriteString(",")
		}
	}

	webhook.Name = name
	webhook.Url = url
	webhook.RequestType = requestType
	webhook.Events = eventString.String()

	err = ct.services.Webhook.Save(webhook)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/webhooks")
		return
	}

	c.Redirect(http.StatusSeeOther, "/webhooks")
}

func (ct *Controller) deleteWebhook(c *gin.Context) {
	id, ok := utils.GetParamOrRedirect(c, "/webhooks", "id")
	if !ok {
		return
	}

	ct.services.Webhook.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/user")
}
