package utils

import (
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/embedfs"
	"github.com/janmarkuslanger/nuricms/internal/server"
)

var templatesFS fs.FS = embedfs.TemplatesFS

func SetTemplatesFS(f fs.FS) {
	templatesFS = f
}

var funcMap = template.FuncMap{
	"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	"eq":       func(a, b any) bool { return a == b },
	"add":      func(a, b int) int { return a + b },
	"sub":      func(a, b int) int { return a - b },
	"in":       func(s string, list []string) bool { return slices.Contains(list, s) },
	"split":    strings.Split,
}

var RenderWithLayoutHTTP = func(ctx server.Context, contentTemplate string, data map[string]any, statusCode int) {
	_, err := ctx.Request.Cookie("auth_token")
	data["IsLoggedIn"] = err == nil

	tmpl, err := template.New("layout.tmpl").
		Funcs(funcMap).
		ParseFS(
			templatesFS,
			"templates/base/layout.tmpl",
			"templates/"+contentTemplate,
		)

	if err != nil {
		http.Error(ctx.Writer, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Writer.WriteHeader(statusCode)
	if err := tmpl.ExecuteTemplate(ctx.Writer, "layout.tmpl", data); err != nil {
		http.Error(ctx.Writer, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

var RenderWithLayout = func(c *gin.Context, contentTemplate string, data gin.H, statusCode int) {
	if _, ok := c.Get("userID"); ok {
		data["IsLoggedIn"] = true
	} else {
		data["IsLoggedIn"] = false
	}

	tmpl, err := template.New("layout.tmpl").
		Funcs(funcMap).
		ParseFS(
			templatesFS,
			"templates/base/layout.tmpl",
			"templates/"+contentTemplate,
		)

	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	c.Status(statusCode)

	if err := tmpl.ExecuteTemplate(c.Writer, "layout.tmpl", data); err != nil {
		c.String(http.StatusInternalServerError, "Template error: %v", err)
	}
}
