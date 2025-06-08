package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 定义错误码常量
const (
	// 成功
	Success ErrorCode = 0

	// 通用错误 (1000-1999)
	InternalError   ErrorCode = 1000
	InvalidParams   ErrorCode = 1001
	NotFound        ErrorCode = 1002
	AlreadyExists   ErrorCode = 1003
	Forbidden       ErrorCode = 1004
	TooManyRequests ErrorCode = 1005

	// 认证错误 (2000-2999)
	Unauthorized  ErrorCode = 2000
	InvalidToken  ErrorCode = 2001
	TokenExpired  ErrorCode = 2002
	InvalidAPIKey ErrorCode = 2003
	LoginFailed   ErrorCode = 2004

	// 业务错误 (3000-3999)
	ProjectNotFound     ErrorCode = 3000
	ProjectExists       ErrorCode = 3001
	LanguageNotFound    ErrorCode = 3002
	LanguageExists      ErrorCode = 3003
	TranslationNotFound ErrorCode = 3004
	UserNotFound        ErrorCode = 3005
	InvalidFileFormat   ErrorCode = 3006

	// 数据库错误 (4000-4999)
	DatabaseError    ErrorCode = 4000
	TransactionError ErrorCode = 4001
	ConnectionError  ErrorCode = 4002

	// 外部服务错误 (5000-5999)
	ExternalAPIError ErrorCode = 5000
)

// 错误码对应的HTTP状态码映射
var httpStatusMap = map[ErrorCode]int{
	Success:             http.StatusOK,
	InternalError:       http.StatusInternalServerError,
	InvalidParams:       http.StatusBadRequest,
	NotFound:            http.StatusNotFound,
	AlreadyExists:       http.StatusConflict,
	Forbidden:           http.StatusForbidden,
	TooManyRequests:     http.StatusTooManyRequests,
	Unauthorized:        http.StatusUnauthorized,
	InvalidToken:        http.StatusUnauthorized,
	TokenExpired:        http.StatusUnauthorized,
	InvalidAPIKey:       http.StatusUnauthorized,
	LoginFailed:         http.StatusUnauthorized,
	ProjectNotFound:     http.StatusNotFound,
	ProjectExists:       http.StatusConflict,
	LanguageNotFound:    http.StatusNotFound,
	LanguageExists:      http.StatusConflict,
	TranslationNotFound: http.StatusNotFound,
	UserNotFound:        http.StatusNotFound,
	InvalidFileFormat:   http.StatusBadRequest,
	DatabaseError:       http.StatusInternalServerError,
	TransactionError:    http.StatusInternalServerError,
	ConnectionError:     http.StatusServiceUnavailable,
	ExternalAPIError:    http.StatusBadGateway,
}

// 错误码对应的消息映射
var messageMap = map[ErrorCode]string{
	Success:             "成功",
	InternalError:       "内部服务器错误",
	InvalidParams:       "参数错误",
	NotFound:            "资源未找到",
	AlreadyExists:       "资源已存在",
	Forbidden:           "权限不足",
	TooManyRequests:     "请求过于频繁",
	Unauthorized:        "未授权访问",
	InvalidToken:        "无效的令牌",
	TokenExpired:        "令牌已过期",
	InvalidAPIKey:       "无效的API密钥",
	LoginFailed:         "登录失败",
	ProjectNotFound:     "项目不存在",
	ProjectExists:       "项目已存在",
	LanguageNotFound:    "语言不存在",
	LanguageExists:      "语言已存在",
	TranslationNotFound: "翻译不存在",
	UserNotFound:        "用户不存在",
	InvalidFileFormat:   "文件格式错误",
	DatabaseError:       "数据库操作失败",
	TransactionError:    "事务操作失败",
	ConnectionError:     "连接失败",
	ExternalAPIError:    "外部服务错误",
}

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Err     error     `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// HTTPStatus 获取对应的HTTP状态码
func (e *AppError) HTTPStatus() int {
	if status, ok := httpStatusMap[e.Code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithError 添加原始错误
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	if e.Details == "" && err != nil {
		e.Details = err.Error()
	}
	return e
}

// NewAppError 创建新的应用错误
func NewAppError(code ErrorCode, details ...string) *AppError {
	message := messageMap[code]
	if message == "" {
		message = "未知错误"
	}

	err := &AppError{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		err.Details = details[0]
	}

	return err
}

// NewAppErrorWithError 创建包含原始错误的应用错误
func NewAppErrorWithError(code ErrorCode, originalErr error, details ...string) *AppError {
	err := NewAppError(code, details...)
	return err.WithError(originalErr)
}

// 便捷函数，创建常用错误
func NewInternalError(details ...string) *AppError {
	return NewAppError(InternalError, details...)
}

func NewInvalidParamsError(details ...string) *AppError {
	return NewAppError(InvalidParams, details...)
}

func NewNotFoundError(details ...string) *AppError {
	return NewAppError(NotFound, details...)
}

func NewUnauthorizedError(details ...string) *AppError {
	return NewAppError(Unauthorized, details...)
}

func NewForbiddenError(details ...string) *AppError {
	return NewAppError(Forbidden, details...)
}

func NewAlreadyExistsError(details ...string) *AppError {
	return NewAppError(AlreadyExists, details...)
}

// 业务错误便捷函数
func NewProjectNotFoundError() *AppError {
	return NewAppError(ProjectNotFound)
}

func NewProjectExistsError() *AppError {
	return NewAppError(ProjectExists)
}

func NewLanguageNotFoundError() *AppError {
	return NewAppError(LanguageNotFound)
}

func NewLanguageExistsError() *AppError {
	return NewAppError(LanguageExists)
}

func NewTranslationNotFoundError() *AppError {
	return NewAppError(TranslationNotFound)
}

func NewUserNotFoundError() *AppError {
	return NewAppError(UserNotFound)
}

func NewDatabaseError(originalErr error) *AppError {
	return NewAppErrorWithError(DatabaseError, originalErr)
}

func NewInvalidToken(details ...string) *AppError {
	return NewAppError(InvalidToken, details...)
}
