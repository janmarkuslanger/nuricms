package content

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseCollectionID(c *gin.Context) (uint, bool) {
	var collectionID uint

	collectionParam := c.Param("id")
	collectionID64, err := strconv.ParseUint(collectionParam, 10, 0)
	if err != nil {
		return collectionID, false
	}

	return uint(collectionID64), true
}
