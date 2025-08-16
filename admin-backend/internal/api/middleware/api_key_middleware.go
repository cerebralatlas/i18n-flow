package middleware

import (
	"i18n-flow/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuthMiddleware 验证API Key中间件
func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		// 获取配置
		cfg, err := config.Load()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
			c.Abort()
			return
		}

		// 如果API Key为空或不匹配
		if apiKey == "" || apiKey != cfg.CLI.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
