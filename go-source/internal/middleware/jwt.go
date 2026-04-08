package middleware

import (
	"go-source/internal/response"
	"go-source/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "token无效或已过期")
			c.Abort()
			return
		}
		c.Set("adminId", claims.AdminID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
