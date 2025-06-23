package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParsePagination(c *gin.Context) (page int, pageSize int) {
	pageParam := c.DefaultQuery("page", "1")
	pageSizeParam := c.DefaultQuery("pageSize", "10")

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
