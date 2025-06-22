package content

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/config"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.FieldTemplates = map[string]string{
		string(model.FieldTypeText): "text_test",
	}
}

func TestRenderFields_RendersValidTemplate(t *testing.T) {
	field := model.Field{
		FieldType: model.FieldTypeText,
		Alias:     "title",
	}

	fields := []FieldContent{
		{
			Field: field,
			Values: []model.ContentValue{
				{Value: "Hello", Field: field},
			},
		},
	}

	html := RenderFields(fields)
	assert.Len(t, html, 1)
	assert.Contains(t, string(html[0]), "Hello")
}

func TestRenderFields_InvalidTemplate_ReturnsEmpty(t *testing.T) {
	field := model.Field{
		FieldType: "nonexistent",
		Alias:     "broken",
	}

	fields := []FieldContent{
		{
			Field: field,
			Values: []model.ContentValue{
				{Value: "N/A", Field: field},
			},
		},
	}

	html := RenderFields(fields)
	assert.Len(t, html, 0)
}

func TestContentToFieldContent_AssignsValuesCorrectly(t *testing.T) {
	field := model.Field{Alias: "title"}
	value := model.ContentValue{Field: field, Value: "Test"}

	collection := model.Collection{
		Fields: []model.Field{field},
	}

	content := model.Content{
		ContentValues: []model.ContentValue{value},
	}

	fields := ContentToFieldContent(content, collection, nil, nil)
	assert.Len(t, fields, 1)
	assert.Equal(t, "Test", fields["title"].Values[0].Value)
}

func TestRenderFieldsByContent_RendersTemplate(t *testing.T) {
	field := model.Field{
		Alias:     "title",
		FieldType: model.FieldTypeText,
	}
	value := model.ContentValue{
		Field: field,
		Value: "Rendered",
	}

	collection := model.Collection{
		Fields: []model.Field{field},
	}

	content := model.Content{
		ContentValues: []model.ContentValue{value},
	}

	html := RenderFieldsByContent(content, collection, nil, nil)
	assert.Len(t, html, 1)
	assert.Contains(t, string(html[0]), "Rendered")
}

func TestContentsToContentGroup_GroupingWorks(t *testing.T) {
	field := model.Field{Alias: "title"}
	value := model.ContentValue{Field: field, Value: "One"}
	content := model.Content{ContentValues: []model.ContentValue{value}}

	groups := ContentsToContentGroup([]model.Content{content})
	assert.Len(t, groups, 1)
	assert.Equal(t, "One", groups[0].ValuesByField["title"][0].Value)
}
