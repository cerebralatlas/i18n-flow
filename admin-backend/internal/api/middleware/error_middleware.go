package middleware

import (
	"fmt"
	"i18n-flow/utils"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 获取请求信息
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("request_id", getRequestIDFromContext(c)),
		}

		// 添加用户信息（如果存在）
		if userID, exists := c.Get("userID"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		if err, ok := recovered.(string); ok {
			utils.ErrorLog("Panic recovered", append(fields,
				zap.String("error", err),
				zap.String("stack", string(debug.Stack())),
			)...)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器发生异常"})
		} else if err, ok := recovered.(error); ok {
			utils.ErrorLog("Panic recovered", append(fields,
				zap.Error(err),
				zap.String("stack", string(debug.Stack())),
			)...)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器发生异常"})
		} else {
			utils.ErrorLog("Panic recovered", append(fields,
				zap.Any("error", recovered),
				zap.String("stack", string(debug.Stack())),
			)...)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器发生异常"})
		}
		c.Abort()
	})
}

// RequestValidationMiddleware 请求验证中间件
func RequestValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查Content-Type（对于POST、PUT请求）
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			contentType := c.GetHeader("Content-Type")
			if contentType != "" && contentType != "application/json" && contentType != "multipart/form-data" {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("不支持的Content-Type: %s", contentType)})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// getRequestIDFromContext 从上下文获取请求ID
func getRequestIDFromContext(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// NotFoundHandler 404处理器
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("路由 %s %s 不存在", c.Request.Method, c.Request.URL.Path)})
	}
}
