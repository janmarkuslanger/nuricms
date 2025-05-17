package template

import (
	"bytes"
	"text/template"
)

func RenderTemplate(templatePath string, data any) (string, error) {
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
