package webhook

import (
	"net/http"
	"strings"

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
	secure := r.Group("/webhooks", middleware.Userauth(h.services.User))

	secure.GET("/", middleware.Roleauth(model.RoleAdmin), h.showWebhook)
	secure.GET("/create", middleware.Roleauth(model.RoleAdmin), h.showCreateWebhook)
	secure.POST("/create", middleware.Roleauth(model.RoleAdmin), h.createWebhook)
	secure.GET("/edit/:id", middleware.Roleauth(model.RoleAdmin), h.showEditWebhook)
	secure.POST("/edit/:id", middleware.Roleauth(model.RoleAdmin), h.editWebhook)
	secure.POST("/delete/:id", middleware.Roleauth(model.RoleAdmin), h.deleteWebhook)
}

func (h *Handler) showWebhook(c *gin.Context) {
	webhooks, _ := h.services.Webhook.List()

	utils.RenderWithLayout(c, "webhook/index.tmpl", gin.H{
		"Webhooks": webhooks,
	}, http.StatusOK)
}

func (h *Handler) showCreateWebhook(c *gin.Context) {
	utils.RenderWithLayout(c, "webhook/create_or_edit.tmpl", gin.H{
		"RequestTypes": model.GetRequestTypes(),
		"EventTypes":   model.GetWebhookEvents(),
	}, http.StatusOK)
}

func (h *Handler) createWebhook(c *gin.Context) {
	name := c.PostForm("name")
	url := c.PostForm("url")
	requestType := model.RequestType(c.PostForm("request_type"))

	events := map[model.EventType]bool{}

	for _, event := range model.GetWebhookEvents() {
		isActive := c.PostForm(string(event)) == "on"
		events[event] = isActive
	}

	_, err := h.services.Webhook.Create(name, url, requestType, events)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/webhooks")
		return
	}

	c.Redirect(http.StatusSeeOther, "/webhooks")
}

func (h *Handler) showEditWebhook(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/webhooks")
	}

	webhook, err := h.services.Webhook.FindByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/webhooks")
	}

	utils.RenderWithLayout(c, "webhook/create_or_edit.tmpl", gin.H{
		"RequestTypes":    model.GetRequestTypes(),
		"EventTypes":      model.GetWebhookEvents(),
		"Webhook":         webhook,
		"EventTypeValues": strings.Split(webhook.Events, ","),
	}, http.StatusOK)
}

func (h *Handler) editWebhook(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))
	if !ok {
		c.Redirect(http.StatusSeeOther, "/webhooks")
	}

	webhook, err := h.services.Webhook.FindByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/webhooks")
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

	err = h.services.Webhook.Save(webhook)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/webhooks")
		return
	}

	c.Redirect(http.StatusSeeOther, "/webhooks")
}

func (h *Handler) deleteWebhook(c *gin.Context) {
	id, ok := utils.StringToUint(c.Param("id"))

	if !ok {
		c.Redirect(http.StatusSeeOther, "/user")
	}

	h.services.Webhook.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/user")
}
