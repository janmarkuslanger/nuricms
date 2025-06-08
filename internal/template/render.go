package template

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/janmarkuslanger/nuricms/internal/embedfs"
)

func RenderTemplate(templatePath string, data any) (string, error) {
	tpl := strings.TrimPrefix(templatePath, "/")

	rawTemplate, err := template.ParseFS(embedfs.TemplatesFS, tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := rawTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
