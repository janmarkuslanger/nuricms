package testutils

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

func MakeGETContext(path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", path, nil)
	return c, w
}

func MakePOSTContext(path string, formData gin.H) (*gin.Context, *httptest.ResponseRecorder) {
	body := url.Values{}
	for key, value := range formData {
		body.Set(key, fmt.Sprintf("%v", value))
	}
	req := httptest.NewRequest("POST", path, strings.NewReader(body.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func MakeBrokenPOSTContext(path string) (*gin.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("POST", path, &brokenReader{})
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

type brokenReader struct{}

func (b *brokenReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("broken reader")
}

func SetParam(c *gin.Context, key, value string) {
	c.Params = append(c.Params, gin.Param{Key: key, Value: value})
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
