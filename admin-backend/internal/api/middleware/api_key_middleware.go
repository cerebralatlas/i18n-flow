package middleware

import (
	"i18n-flow/internal/api/response"
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
			response.InternalServerError(c, "Failed to load configuration")
			return
		}

		// 如果API Key为空或不匹配
		if apiKey == "" || apiKey != cfg.CLI.APIKey {
			response.Error(c, http.StatusUnauthorized, "INVALID_API_KEY", "Invalid or missing API key")
			return
		}

		c.Next()
	}
}
