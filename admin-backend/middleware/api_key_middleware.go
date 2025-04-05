package middleware

import (
	"i18n-flow/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuthMiddleware 验证API Key中间件
func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		// 如果API Key为空或不匹配
		if apiKey == "" || apiKey != config.GetConfig().CLI.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
