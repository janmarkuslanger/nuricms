package utils

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RenderTemplate(templateName string, data interface{}) (string, error) {
	templatePath := filepath.Join("templates", templateName+".html")
	rawTemplate, err := template.ParseFiles(templatePath)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := rawTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderWithLayout(c *gin.Context, contentTemplate string, data gin.H, statusCode int) {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"eq": func(a, b interface{}) bool {
			return a == b
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
