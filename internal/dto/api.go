package dto

import (
	"time"

	"github.com/janmarkuslanger/nuricms/internal/model"
)

type ErrorDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type MetaData struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type Pagination struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}

type ApiResponse struct {
	Success    bool         `json:"success"`
	Data       interface{}  `json:"data,omitempty"`
	Error      *ErrorDetail `json:"error,omitempty"`
	Meta       *MetaData    `json:"meta,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
}

type CollectionResponse struct {
	ID    uint   `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Alias string `json:"alias,omitempty"`
}

type AssetResponse struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type ContentItemResponse struct {
	ID         uint               `json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Values     map[string]any     `json:"values"`
	Collection CollectionResponse `json:"collection"`
}

type ContentValueResponse struct {
	ID         uint                `json:"id"`
	Value      any                 `json:"value"`
	FieldType  model.FieldType     `json:"field_type"`
	Collection *CollectionResponse `json:"collection,omitempty"`
	Asset      *AssetResponse      `json:"asset,omitempty"`
}
