package content

import (
	"html/template"
	"path/filepath"

	"github.com/janmarkuslanger/nuricms/config"
	"github.com/janmarkuslanger/nuricms/model"
	pkgtemplate "github.com/janmarkuslanger/nuricms/pkg/template"
)

func RenderFields(fields []model.Field) []template.HTML {
	var htmlFields []template.HTML

	for _, field := range fields {
		templateName := config.FieldTemplates[string(field.FieldType)]
		templatePath := filepath.Join("templates", templateName+".html")
		templateContent, err := pkgtemplate.RenderTemplate(templatePath, field)

		if err == nil {
			safeHtml := template.HTML(templateContent)
			htmlFields = append(htmlFields, safeHtml)
		}
	}

	return htmlFields
}

type FieldContent struct {
	Field  model.Field
	Values []model.ContentValue
}

func CreateFieldContentsByContent(content model.Content) map[string]FieldContent {
	fields := make(map[string]FieldContent)

	for _, contentValue := range content.ContentValues {
		fieldAlias := contentValue.Field.Alias

		v, ok := fields[fieldAlias]

		if ok {
			v.Values = append(v.Values, contentValue)
			continue
		}

		values := []model.ContentValue{contentValue}

		fields[fieldAlias] = FieldContent{
			Field:  contentValue.Field,
			Values: values,
		}
	}

	return fields
}

func RenderFieldsByContent(content model.Content) []template.HTML {
	var htmlFields []template.HTML

	fields := CreateFieldContentsByContent(content)

	for _, field := range fields {
		templateName := config.FieldTemplates[string(field.Field.FieldType)]
		templatePath := filepath.Join("templates", templateName+".html")
		templateContent, err := pkgtemplate.RenderTemplate(templatePath, field)

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

func CreateContentGroupsByField(contents []model.Content) []ContentGroup {
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
