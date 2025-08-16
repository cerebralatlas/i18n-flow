package domain

import "errors"

// 定义领域层错误
var (
	// 用户相关错误
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
	ErrUserExists      = errors.New("用户已存在")
	ErrInvalidToken    = errors.New("无效的令牌")

	// 项目相关错误
	ErrProjectNotFound = errors.New("项目不存在")
	ErrProjectExists   = errors.New("项目已存在")
	ErrInvalidSlug     = errors.New("无效的项目标识")

	// 语言相关错误
	ErrLanguageNotFound = errors.New("语言不存在")
	ErrLanguageExists   = errors.New("语言已存在")
	ErrInvalidLanguage  = errors.New("无效的语言代码")

	// 翻译相关错误
	ErrTranslationNotFound = errors.New("翻译不存在")
	ErrTranslationExists   = errors.New("翻译已存在")
	ErrInvalidKey          = errors.New("无效的翻译键")

	// 通用错误
	ErrInvalidInput  = errors.New("无效的输入参数")
	ErrInternalError = errors.New("内部服务器错误")
	ErrUnauthorized  = errors.New("未授权访问")
	ErrForbidden     = errors.New("禁止访问")
)
