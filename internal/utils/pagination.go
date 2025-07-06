package utils

import (
	"net/http"
	"strconv"
)

func ParsePagination(r *http.Request) (page int, pageSize int) {
	pageParam := DefaultQuery(r, "page", "1")
	pageSizeParam := DefaultQuery(r, "pageSize", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	pageSize, err = strconv.Atoi(pageSizeParam)
	if err != nil {
		pageSize = 10
	}

	return page, pageSize
}
