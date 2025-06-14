package dto

type ContentWithValues struct {
	CollectionID uint
	ContentID    uint
	FormData     map[string][]string
}
