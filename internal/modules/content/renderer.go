package content

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/janmarkuslanger/nuricms/internal/embedfs"
	"github.com/janmarkuslanger/nuricms/internal/globals"
	"github.com/janmarkuslanger/nuricms/internal/model"
	utilstemplate "github.com/janmarkuslanger/nuricms/internal/template"
)

type FieldContent struct {
	Field   model.Field
	Values  []model.ContentValue
	Content []model.Content
	Assets  []model.Asset
}

func renderField(content FieldContent) (template.HTML, error) {
	var html template.HTML
	templateName := globals.FieldTemplates[string(content.Field.FieldType)]
	templatePath := filepath.Join("templates", templateName+".tmpl")
	templateContent, err := utilstemplate.RenderTemplate(embedfs.TemplatesFS, templatePath, content)
	if err != nil {
		return html, err
	}

	return template.HTML(templateContent), nil
}

func RenderFields(fields []FieldContent) []template.HTML {
	var htmlFields []template.HTML

	for _, field := range fields {
		html, err := renderField(field)
		if err == nil {
			htmlFields = append(htmlFields, html)
		}
	}

	return htmlFields
}

type DataContext struct {
	Collection model.Collection
	Contents   []model.Content
	Assets     []model.Asset
}

func ContentToFieldContent(content model.Content, ctx DataContext) map[string]FieldContent {
	fields := make(map[string]FieldContent)

	for _, field := range ctx.Collection.Fields {
		values := make([]model.ContentValue, 0)
		fields[field.Alias] = FieldContent{
			Field:   field,
			Values:  values,
			Content: ctx.Contents,
			Assets:  ctx.Assets,
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

func RenderFieldsByContent(content model.Content, ctx DataContext) []template.HTML {
	var htmlFields []template.HTML

	fields := ContentToFieldContent(content, ctx)

	for _, field := range fields {
		html, err := renderField(field)
		if err != nil {
			fmt.Printf("error rendering field %v: %v\n", field, err)
			continue
		}

		htmlFields = append(htmlFields, html)
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
