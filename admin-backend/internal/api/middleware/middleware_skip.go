package middleware

import (
	"github.com/gin-gonic/gin"
)

// IsSwaggerPath 检查请求路径是否为 Swagger 路径
func IsSwaggerPath(c *gin.Context) bool {
	path := c.Request.URL.Path
	return len(path) >= 8 && path[:8] == "/swagger"
}

// SkipForSwagger 创建一个跳过 Swagger 路径的中间件包装器
func SkipForSwagger(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSwaggerPath(c) {
			c.Next()
			return
		}
		handler(c)
	}
}
