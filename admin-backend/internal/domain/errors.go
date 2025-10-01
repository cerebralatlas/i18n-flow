package domain

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
)

// AppError 应用程序错误
type AppError struct {
	Type      ErrorType              `json:"type"`
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Cause     error                  `json:"-"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Cause
}

// HTTPStatus 返回对应的HTTP状态码
func (e *AppError) HTTPStatus() int {
	switch e.Type {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewAppError 创建新的应用程序错误
func NewAppError(errType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:      errType,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewAppErrorWithCause 创建带原因的应用程序错误
func NewAppErrorWithCause(errType ErrorType, code, message string, cause error) *AppError {
	return &AppError{
		Type:      errType,
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
	}
}

// NewAppErrorWithDetails 创建带详细信息的应用程序错误
func NewAppErrorWithDetails(errType ErrorType, code, message, details string) *AppError {
	return &AppError{
		Type:      errType,
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// NewAppErrorWithContext 创建带上下文的应用程序错误
func NewAppErrorWithContext(errType ErrorType, code, message string, context map[string]interface{}) *AppError {
	return &AppError{
		Type:      errType,
		Code:      code,
		Message:   message,
		Context:   context,
		Timestamp: time.Now(),
	}
}

// 预定义的领域错误
var (
	// 用户相关错误
	ErrUserNotFound      = NewAppError(ErrorTypeNotFound, "USER_NOT_FOUND", "用户不存在")
	ErrInvalidPassword   = NewAppError(ErrorTypeUnauthorized, "INVALID_PASSWORD", "密码错误")
	ErrUserExists        = NewAppError(ErrorTypeConflict, "USER_EXISTS", "用户已存在")
	ErrEmailExists       = NewAppError(ErrorTypeConflict, "EMAIL_EXISTS", "邮箱已存在")
	ErrInvalidToken      = NewAppError(ErrorTypeUnauthorized, "INVALID_TOKEN", "无效的令牌")
	ErrInvalidRole       = NewAppError(ErrorTypeValidation, "INVALID_ROLE", "无效的角色")
	ErrCannotDeleteAdmin = NewAppError(ErrorTypeForbidden, "CANNOT_DELETE_ADMIN", "不能删除管理员用户")

	// 项目相关错误
	ErrProjectNotFound = NewAppError(ErrorTypeNotFound, "PROJECT_NOT_FOUND", "项目不存在")
	ErrProjectExists   = NewAppError(ErrorTypeConflict, "PROJECT_EXISTS", "项目已存在")
	ErrInvalidSlug     = NewAppError(ErrorTypeValidation, "INVALID_SLUG", "无效的项目标识")

	// 语言相关错误
	ErrLanguageNotFound = NewAppError(ErrorTypeNotFound, "LANGUAGE_NOT_FOUND", "语言不存在")
	ErrLanguageExists   = NewAppError(ErrorTypeConflict, "LANGUAGE_EXISTS", "语言已存在")
	ErrInvalidLanguage  = NewAppError(ErrorTypeValidation, "INVALID_LANGUAGE", "无效的语言代码")

	// 翻译相关错误
	ErrTranslationNotFound = NewAppError(ErrorTypeNotFound, "TRANSLATION_NOT_FOUND", "翻译不存在")
	ErrTranslationExists   = NewAppError(ErrorTypeConflict, "TRANSLATION_EXISTS", "翻译已存在")
	ErrInvalidKey          = NewAppError(ErrorTypeValidation, "INVALID_KEY", "无效的翻译键")

	// 项目成员相关错误
	ErrMemberNotFound    = NewAppError(ErrorTypeNotFound, "MEMBER_NOT_FOUND", "项目成员不存在")
	ErrMemberExists      = NewAppError(ErrorTypeConflict, "MEMBER_EXISTS", "用户已是项目成员")
	ErrInsufficientPerm  = NewAppError(ErrorTypeForbidden, "INSUFFICIENT_PERMISSION", "权限不足")
	ErrCannotRemoveOwner = NewAppError(ErrorTypeForbidden, "CANNOT_REMOVE_OWNER", "不能移除项目所有者")

	// 通用错误
	ErrInvalidInput  = NewAppError(ErrorTypeValidation, "INVALID_INPUT", "无效的输入参数")
	ErrInternalError = NewAppError(ErrorTypeInternal, "INTERNAL_ERROR", "内部服务器错误")
	ErrUnauthorized  = NewAppError(ErrorTypeUnauthorized, "UNAUTHORIZED", "未授权访问")
	ErrForbidden     = NewAppError(ErrorTypeForbidden, "FORBIDDEN", "禁止访问")
)

// IsAppError 检查是否为应用程序错误
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// WrapError 包装普通错误为应用程序错误
func WrapError(err error, errType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:      errType,
		Code:      code,
		Message:   message,
		Cause:     err,
		Timestamp: time.Now(),
	}
}
