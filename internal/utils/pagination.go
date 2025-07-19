package utils

import (
	"net/http"
	"strconv"
)

func ParsePagination(r *http.Request) (page int, pageSize int) {
	page, err := strconv.Atoi(DefaultQuery(r, "page", "1"))
	if err != nil {
		page = 1
	}

	pageSize, err = strconv.Atoi(DefaultQuery(r, "pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	return page, pageSize
}

func CalcTotalPages(total int64, pageSize int) int64 {
	return (total + int64(pageSize) - 1) / int64(pageSize)
}
