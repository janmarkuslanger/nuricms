package utils

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RenderWithLayout(c *gin.Context, contentTemplate string, data gin.H, statusCode int) {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"add": func(a int, b int) int {
			return a + b
		},
	}

	tmpl := template.Must(template.New("layout.html").
		Funcs(funcMap).
		ParseFiles(
			"templates/base/layout.html",
			"templates/"+contentTemplate,
		),
	)

	c.Status(statusCode)
	if err := tmpl.ExecuteTemplate(c.Writer, "layout.html", data); err != nil {
		c.String(http.StatusInternalServerError, "Template error: %v", err)
	}
}
