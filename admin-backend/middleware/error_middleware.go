package middleware

import (
	"fmt"
	"i18n-flow/errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s\n%s", err, debug.Stack())
			errors.ErrorResponse(c, errors.NewInternalError("服务器发生异常"))
		} else if err, ok := recovered.(error); ok {
			log.Printf("Panic recovered: %s\n%s", err.Error(), debug.Stack())
			errors.ErrorResponse(c, errors.NewInternalError("服务器发生异常"))
		} else {
			log.Printf("Panic recovered: %v\n%s", recovered, debug.Stack())
			errors.ErrorResponse(c, errors.NewInternalError("服务器发生异常"))
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
				errors.BadRequestResponse(c, fmt.Sprintf("不支持的Content-Type: %s", contentType))
				return
			}
		}

		c.Next()
	}
}

// NotFoundHandler 404处理器
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		errors.NotFoundResponse(c, fmt.Sprintf("路由 %s %s 不存在", c.Request.Method, c.Request.URL.Path))
	}
}
