package handler

import (
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type ListHandler[T any] interface {
	List(page, pageSize int) ([]T, int64, error)
}

func HandleList[T any](ctx server.Context, s ListHandler[T], template string) {
	page, pageSize := utils.ParsePagination(ctx.Request)

	items, totalCount, _ := s.List(page, pageSize)

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	utils.RenderWithLayoutHTTP(ctx, template, map[string]any{
		"Items":       items,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"CurrentPage": page,
		"PageSize":    pageSize,
	}, http.StatusOK)
}
