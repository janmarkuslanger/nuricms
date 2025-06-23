package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetParamOrRedirect(c *gin.Context, r string, p string) (uint, bool) {
	id, ok := StringToUint(c.Param(p))
	if !ok {
		c.Redirect(http.StatusSeeOther, r)
	}

	return id, ok
}
