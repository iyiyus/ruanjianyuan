package middleware

import (
	"go-source/internal/handler"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireInstalled 未安装时拦截所有请求跳转引导页
func RequireInstalled() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 放行安装相关路径
		if path == "/install" ||
			strings.HasPrefix(path, "/install/") ||
			strings.HasPrefix(path, "/api/install/") {
			c.Next()
			return
		}
		if !handler.IsInstalled() {
			c.Redirect(http.StatusFound, "/install?step=0")
			c.Abort()
			return
		}
		c.Next()
	}
}
