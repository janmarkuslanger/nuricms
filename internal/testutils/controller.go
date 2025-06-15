package testutils

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

func MakeGETContext(path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", path, nil)
	return c, w
}

func StubRenderWithLayout() func() {
	orig := utils.RenderWithLayout
	utils.RenderWithLayout = func(c *gin.Context, contentTemplate string, data gin.H, statusCode int) {
		c.Status(statusCode)
		c.Writer.WriteString("RENDERED:" + contentTemplate)
	}
	return func() {
		utils.RenderWithLayout = orig
	}
}
