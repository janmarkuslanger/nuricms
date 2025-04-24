package content

import (
	"html/template"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

func RenderFields(fields []model.Field) []template.HTML {

	var htmlFields []template.HTML

	for _, field := range fields {
		var templatePath string

		if field.FieldType == "Text" {
			templatePath = "template_fields/text"
		}

		templateContent, err := utils.RenderTemplate(templatePath, field)

		if err == nil {
			safeHtml := template.HTML(templateContent)
			htmlFields = append(htmlFields, safeHtml)
		}

	}

	return htmlFields
}
