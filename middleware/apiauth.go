package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/service"
)

func ApikeyAuth(keySvc *service.ApikeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-API-Key")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing X-API-Key header",
			})
			return
		}

		if err := keySvc.Validate(token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Next()
	}
}
