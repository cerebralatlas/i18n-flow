package errors

import (
	"github.com/gin-gonic/gin"
)

// Response 标准响应结构
type Response struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Details string      `json:"details,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    Success,
		Message: messageMap[Success],
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, err error) {
	var appErr *AppError

	// 判断是否为AppError类型
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		// 如果不是AppError，包装为内部错误
		appErr = NewInternalError(err.Error())
	}

	c.JSON(appErr.HTTPStatus(), Response{
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: appErr.Details,
	})
}

// AbortWithError 中断请求并返回错误
func AbortWithError(c *gin.Context, err error) {
	ErrorResponse(c, err)
	c.Abort()
}

// 便捷响应函数
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(200, Response{
		Code:    Success,
		Message: message,
		Data:    data,
	})
}

func BadRequestResponse(c *gin.Context, details string) {
	ErrorResponse(c, NewInvalidParamsError(details))
}

func NotFoundResponse(c *gin.Context, details string) {
	ErrorResponse(c, NewNotFoundError(details))
}

func UnauthorizedResponse(c *gin.Context, details string) {
	ErrorResponse(c, NewUnauthorizedError(details))
}

func ForbiddenResponse(c *gin.Context, details string) {
	ErrorResponse(c, NewForbiddenError(details))
}

func InternalServerErrorResponse(c *gin.Context, details string) {
	ErrorResponse(c, NewInternalError(details))
}
