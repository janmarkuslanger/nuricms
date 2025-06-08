package dto

import "time"

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
