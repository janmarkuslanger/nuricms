package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/service"
)

func Userauth(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil {
			hdr := c.GetHeader("Authorization")
			if !strings.HasPrefix(hdr, "Bearer ") {
				c.Redirect(http.StatusSeeOther, "/login")
				return
			}
			token = strings.TrimPrefix(hdr, "Bearer ")
		}

		uid, email, role, err := userService.ValidateJWT(token)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		c.Set("userID", uid)
		c.Set("userEmail", email)
		c.Set("userRole", role)

		c.Next()
	}
}

func Roleauth(allowed ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("userRole")
		userRole := val.(model.Role)

		for _, want := range allowed {
			if userRole == want {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
