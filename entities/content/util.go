package content

import (
	"html/template"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/utils"
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

type ContentGroup struct {
	Content       model.Content
	ValuesByField map[string][]model.ContentValue
}

func GroupContentValuesByField(contents []model.Content) []ContentGroup {
	var groups []ContentGroup

	for _, content := range contents {
		fieldGroups := make(map[string][]model.ContentValue)

		for _, value := range content.ContentValues {
			fieldGroups[value.Field.Alias] = append(fieldGroups[value.Field.Alias], value)
		}

		groups = append(groups, ContentGroup{
			Content:       content,
			ValuesByField: fieldGroups,
		})
	}

	return groups
}
