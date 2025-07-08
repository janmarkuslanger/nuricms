package template

import (
	"bytes"
	"embed"
	"html/template"
	"strings"
)

var RenderTemplate = func(embededFs embed.FS, templatePath string, data any) (string, error) {
	tpl := strings.TrimPrefix(templatePath, "/")

	rawTemplate, err := template.ParseFS(embededFs, tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := rawTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
