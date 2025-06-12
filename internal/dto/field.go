package dto

type FieldData struct {
	Name         string
	Alias        string
	CollectionID string
	FieldType    string
	IsList       string
	IsRequired   string
	DisplayField string
}

// Name:         name,
// 		Alias:        alias,
// 		CollectionID: uint(collectionID),
// 		FieldType:    model.FieldType(c.PostForm("field_type")),
// 		IsList:       c.PostForm("is_list") == "on",
// 		IsRequired:   c.PostForm("is_required") == "on",
// 		DisplayField: c.PostForm("display_field") == "on",
