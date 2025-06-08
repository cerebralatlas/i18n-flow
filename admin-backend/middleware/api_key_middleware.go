package middleware

import (
	"i18n-flow/config"
	"i18n-flow/errors"

	"github.com/gin-gonic/gin"
)

// APIKeyAuthMiddleware 验证API Key中间件
func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		// 如果API Key为空或不匹配
		if apiKey == "" || apiKey != config.GetConfig().CLI.APIKey {
			errors.AbortWithError(c, errors.NewAppError(errors.InvalidAPIKey, "Invalid or missing API key"))
			return
		}

		c.Next()
	}
}
