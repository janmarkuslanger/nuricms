package content

import (
	"embed"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	utilstemplate "github.com/janmarkuslanger/nuricms/internal/template"
	"github.com/stretchr/testify/assert"
)

func TestContentToFieldContent(t *testing.T) {
	field := model.Field{Alias: "title", FieldType: "text"}
	value := model.ContentValue{Field: field, Value: "Hello"}

	ctx := DataContext{
		Collection: model.Collection{
			Fields: []model.Field{field},
		},
		Contents: []model.Content{},
		Assets:   []model.Asset{},
	}

	c := model.Content{
		ContentValues: []model.ContentValue{value},
	}

	result := ContentToFieldContent(c, ctx)

	assert.Len(t, result, 1)
	assert.Equal(t, "title", result["title"].Field.Alias)
	assert.Len(t, result["title"].Values, 1)
	assert.Equal(t, "Hello", result["title"].Values[0].Value)
}

func TestContentsToContentGroup(t *testing.T) {
	field := model.Field{Alias: "desc", FieldType: "text"}
	value := model.ContentValue{Field: field, Value: "Content A"}

	c := model.Content{
		ContentValues: []model.ContentValue{value},
	}

	groups := ContentsToContentGroup([]model.Content{c})

	assert.Len(t, groups, 1)
	assert.Equal(t, c, groups[0].Content)
	assert.Len(t, groups[0].ValuesByField["desc"], 1)
	assert.Equal(t, "Content A", groups[0].ValuesByField["desc"][0].Value)
}

func TestRenderFields_Success(t *testing.T) {
	original := utilstemplate.RenderTemplate
	defer func() { utilstemplate.RenderTemplate = original }()

	utilstemplate.RenderTemplate = func(embededFs embed.FS, templatePath string, data any) (string, error) {
		return "<div>Mocked</div>", nil
	}
	field := model.Field{Alias: "body", FieldType: "text"}
	contentField := FieldContent{
		Field:   field,
		Values:  []model.ContentValue{},
		Content: []model.Content{},
		Assets:  []model.Asset{},
	}

	f := RenderFields([]FieldContent{contentField})
	assert.Len(t, f, 1)
}

func TestRenderFieldsByContent(t *testing.T) {
	original := utilstemplate.RenderTemplate
	defer func() { utilstemplate.RenderTemplate = original }()

	utilstemplate.RenderTemplate = func(embededFs embed.FS, templatePath string, data any) (string, error) {
		return "<div>Mocked</div>", nil
	}

	field := model.Field{Alias: "header", FieldType: "text"}
	value := model.ContentValue{Field: field, Value: "Welcome"}

	c := model.Content{
		ContentValues: []model.ContentValue{value},
	}

	ctx := DataContext{
		Collection: model.Collection{
			Fields: []model.Field{field},
		},
		Contents: []model.Content{},
		Assets:   []model.Asset{},
	}

	htmlFields := RenderFieldsByContent(c, ctx)
	assert.NotNil(t, htmlFields)
	assert.Greater(t, len(htmlFields), 0)
}
