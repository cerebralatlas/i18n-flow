package middleware

import (
	"bytes"
	"i18n-flow/utils"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseWriter 包装响应写入器以获取状态码和响应大小
type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
	size       int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.body.Write(b)
	w.size += size
	return size, err
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware 请求日志中间件
func LoggingMiddleware() gin.HandlerFunc {
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
		utils.AccessLog(message, fields...)

		// 错误请求额外记录到错误日志
		if rw.statusCode >= 400 {
			utils.ErrorLog("HTTP Error", fields...)
		}
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单的请求ID生成，实际项目中可以使用UUID
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// shouldLogRequestBody 判断是否应该记录请求体
func shouldLogRequestBody(path string) bool {
	// 不记录敏感路径的请求体
	sensitiveEndpoints := []string{
		"/api/login",
		"/api/refresh",
		"/api/register",
	}

	for _, endpoint := range sensitiveEndpoints {
		if path == endpoint {
			return false
		}
	}

	return true
}

// SkipLoggingMiddleware 跳过日志的中间件（用于健康检查等）
func SkipLoggingMiddleware(paths ...string) gin.HandlerFunc {
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
