package middleware

import (
	"bytes"
	"i18n-flow/internal/utils"
	app_utils "i18n-flow/utils"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EnhancedLoggingMiddleware 增强的日志中间件，集成监控功能
func EnhancedLoggingMiddleware(monitor *utils.SimpleMonitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 包装响应写入器
		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
			statusCode:     200,
		}
		c.Writer = rw

		// 读取请求体（如果需要）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(start)

		// 记录监控指标
		monitor.RecordRequest()

		// 记录慢请求
		if duration > time.Second {
			monitor.RecordSlowRequest()
		}

		// 记录错误请求
		if rw.statusCode >= 400 {
			monitor.RecordError()
		}

		// 收集日志字段
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status_code", rw.statusCode),
			zap.Int("response_size", rw.size),
			zap.Duration("duration", duration),
			zap.String("request_id", getRequestID(c)),
		}

		// 添加用户信息（如果存在）
		if userID, exists := c.Get("userID"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}
		if username, exists := c.Get("username"); exists {
			fields = append(fields, zap.String("username", username.(string)))
		}

		// 添加请求体（仅在debug级别且非敏感路径）
		if shouldLogRequestBody(c.Request.URL.Path) && len(requestBody) > 0 && len(requestBody) < 1024 {
			fields = append(fields, zap.String("request_body", string(requestBody)))
		}

		// 跳过日志记录的路径
		if _, skip := c.Get("skip_logging"); skip {
			return
		}

		// 记录到访问日志
		message := "HTTP Request"
		app_utils.AccessLog(message, fields...)

		// 慢请求警告日志
		if duration > time.Second {
			app_utils.AppWarn("Slow request detected",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
				zap.Duration("duration", duration),
				zap.Int("status", rw.statusCode),
			)
		}

		// 错误请求记录到错误日志
		if rw.statusCode >= 400 {
			if rw.statusCode >= 500 {
				// 5xx 错误使用 AppError 记录更详细信息
				app_utils.AppError("Server error occurred", fields...)
			} else {
				// 4xx 错误使用普通错误日志
				app_utils.ErrorLog("Client error", fields...)
			}
		}
	}
}

// MonitoringStatsMiddleware 监控统计中间件（轻量级版本）
func MonitoringStatsMiddleware(monitor *utils.SimpleMonitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 简单的统计收集
		duration := time.Since(start)
		status := c.Writer.Status()

		// 记录基础指标
		monitor.RecordRequest()

		if duration > time.Second {
			monitor.RecordSlowRequest()
		}

		if status >= 400 {
			monitor.RecordError()
		}

		// 设置响应头，方便调试
		c.Header("X-Response-Time", duration.String())
	}
}

// HealthCheckSkipMiddleware 跳过健康检查的日志记录
func HealthCheckSkipMiddleware(paths ...string) gin.HandlerFunc {
	skipPaths := make(map[string]bool)
	for _, path := range paths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		if skipPaths[c.Request.URL.Path] {
			c.Set("skip_logging", true)
		}
		c.Next()
	}
}
