package utils

import (
	"html/template"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func RenderWithLayout(c *gin.Context, contentTemplate string, data gin.H, statusCode int) {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"eq": func(a, b any) bool {
			return a == b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"in": func(s string, list []string) bool {
			return slices.Contains(list, s)
		},
	}

	if _, ok := c.Get("userID"); ok {
		data["IsLoggedIn"] = true
	} else {
		data["IsLoggedIn"] = false
	}

	tmpl := template.Must(template.New("layout.tmpl").
		Funcs(funcMap).
		ParseFiles(
			"templates/base/layout.tmpl",
			"templates/"+contentTemplate,
		),
	)

	c.Status(statusCode)
	if err := tmpl.ExecuteTemplate(c.Writer, "layout.tmpl", data); err != nil {
		c.String(http.StatusInternalServerError, "Template error: %v", err)
	}
}
