package utils

import (
	"i18n-flow/errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateRequired 验证必填字段
func ValidateRequired(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.NewInvalidParamsError(fieldName + "不能为空")
	}
	return nil
}

// ValidateID 验证ID参数
func ValidateID(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, errors.NewInvalidParamsError(paramName + "参数不能为空")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		return 0, errors.NewInvalidParamsError(paramName + "参数格式错误")
	}

	return uint(id), nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return nil // 允许为空，如果需要必填，另外检查
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return errors.NewInternalError("邮箱格式验证失败")
	}

	if !matched {
		return errors.NewInvalidParamsError("邮箱格式不正确")
	}

	return nil
}

// ValidateStringLength 验证字符串长度
func ValidateStringLength(value, fieldName string, minLen, maxLen int) error {
	length := len(strings.TrimSpace(value))

	if length < minLen {
		return errors.NewInvalidParamsError(fieldName + "长度不能少于" + strconv.Itoa(minLen) + "个字符")
	}

	if maxLen > 0 && length > maxLen {
		return errors.NewInvalidParamsError(fieldName + "长度不能超过" + strconv.Itoa(maxLen) + "个字符")
	}

	return nil
}

// ValidateSlug 验证slug格式（用于URL友好的标识符）
func ValidateSlug(slug string) error {
	if slug == "" {
		return errors.NewInvalidParamsError("slug不能为空")
	}

	// slug只能包含小写字母、数字和连字符，且不能以连字符开头或结尾
	pattern := `^[a-z0-9]+(-[a-z0-9]+)*$`
	matched, err := regexp.MatchString(pattern, slug)
	if err != nil {
		return errors.NewInternalError("slug格式验证失败")
	}

	if !matched {
		return errors.NewInvalidParamsError("slug格式不正确，只能包含小写字母、数字和连字符")
	}

	return nil
}

// ValidateLanguageCode 验证语言代码格式
func ValidateLanguageCode(code string) error {
	if code == "" {
		return errors.NewInvalidParamsError("语言代码不能为空")
	}

	// 支持格式: en, zh, zh-CN, en-US等
	pattern := `^[a-z]{2}(-[A-Z]{2})?$`
	matched, err := regexp.MatchString(pattern, code)
	if err != nil {
		return errors.NewInternalError("语言代码格式验证失败")
	}

	if !matched {
		return errors.NewInvalidParamsError("语言代码格式不正确，应为类似'en'或'zh-CN'的格式")
	}

	return nil
}

// ValidatePagination 验证分页参数
func ValidatePagination(c *gin.Context) (page, pageSize int, err error) {
	// 默认值
	page = 1
	pageSize = 10

	// 解析page参数
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, errors.NewInvalidParamsError("page参数必须为正整数")
		}
	}

	// 解析pageSize参数
	if pageSizeStr := c.DefaultQuery("pageSize", "10"); pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			return 0, 0, errors.NewInvalidParamsError("pageSize参数必须为正整数")
		}

		// 限制最大页面大小
		if pageSize > 100 {
			return 0, 0, errors.NewInvalidParamsError("pageSize参数不能超过100")
		}
	}

	return page, pageSize, nil
}
