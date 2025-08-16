package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta 元数据（用于分页等）
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithStatus 带状态码的成功响应
func SuccessWithStatus(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta 带元数据的成功响应（用于分页）
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Error 错误响应
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
	c.Abort()
}

// ErrorWithDetails 带详细信息的错误响应
func ErrorWithDetails(c *gin.Context, status int, code, message, details string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
	c.Abort()
}

// 预定义的错误响应函数
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func BadRequestWithDetails(c *gin.Context, message, details string) {
	ErrorWithDetails(c, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "CONFLICT", message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "VALIDATION_ERROR", message)
}

func TokenExpired(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "TOKEN_EXPIRED", message)
}

func InvalidToken(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "INVALID_TOKEN", message)
}

// NoContent 无内容响应（通常用于DELETE操作）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}
