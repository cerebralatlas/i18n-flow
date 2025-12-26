package middleware

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装响应写入器以获取状态码和响应大小
type ResponseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
	size       int
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.body.Write(b)
	w.size += size
	return size, err
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// GetRequestID 获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if s, ok := requestID.(string); ok {
			return s
		}
	}
	return ""
}

// GenerateRequestID 生成请求ID（使用加密安全的随机数）
func GenerateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return time.Now().Format("20060102150405") + "-" + hex.EncodeToString(b)
}

// ShouldLogRequestBody 判断是否应该记录请求体
func ShouldLogRequestBody(path string) bool {
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
