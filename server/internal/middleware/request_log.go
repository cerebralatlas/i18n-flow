package middleware

import (
	"time"

	"i18n-flow/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RequestLog 打印请求日志
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("request completed",
			logger.String("path", path),
			logger.String("method", method),
			logger.Int("status", status),
			logger.String("latency", latency.String()),
		)
	}
}
