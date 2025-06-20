package content

import (
	"html/template"
	"path/filepath"

	"github.com/janmarkuslanger/nuricms/internal/config"
	"github.com/janmarkuslanger/nuricms/internal/model"
	utilstemplate "github.com/janmarkuslanger/nuricms/internal/template"
)

type FieldContent struct {
	Field   model.Field
	Values  []model.ContentValue
	Content []model.Content
	Assets  []model.Asset
}

func RenderFields(fields []FieldContent) []template.HTML {
	var htmlFields []template.HTML

	for _, field := range fields {
		templateName := config.FieldTemplates[string(field.Field.FieldType)]
		templatePath := filepath.Join("templates", templateName+".tmpl")
		templateContent, err := utilstemplate.RenderTemplate(templatePath, field)

		if err == nil {
			safeHtml := template.HTML(templateContent)
			htmlFields = append(htmlFields, safeHtml)
		}
	}

	return htmlFields
}

// TODO: too much args
func ContentToFieldContent(content model.Content, collection model.Collection, contents []model.Content, assets []model.Asset) map[string]FieldContent {
	fields := make(map[string]FieldContent)

	for _, field := range collection.Fields {
		values := make([]model.ContentValue, 0)

		fields[field.Alias] = FieldContent{
			Field:   field,
			Values:  values,
			Content: contents,
			Assets:  assets,
		}
	}

	for _, contentValue := range content.ContentValues {
		fieldAlias := contentValue.Field.Alias

		v, ok := fields[fieldAlias]

		if ok {
			v.Values = append(v.Values, contentValue)
			fields[fieldAlias] = v
		}
	}

	return fields
}

// TODO: too much args
func RenderFieldsByContent(content model.Content, collection model.Collection, contents []model.Content, assets []model.Asset) []template.HTML {
	var htmlFields []template.HTML

	fields := ContentToFieldContent(content, collection, contents, assets)

	for _, field := range fields {
		templateName := config.FieldTemplates[string(field.Field.FieldType)]
		templatePath := filepath.Join("templates", templateName+".tmpl")

		templateContent, err := utilstemplate.RenderTemplate(templatePath, field)
		if err != nil {
			continue
		}

		safeHtml := template.HTML(templateContent)
		htmlFields = append(htmlFields, safeHtml)
	}

	return htmlFields
}

type ContentGroup struct {
	Content       model.Content
	ValuesByField map[string][]model.ContentValue
}

func ContentsToContentGroup(contents []model.Content) []ContentGroup {
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
