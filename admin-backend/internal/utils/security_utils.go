package utils

import (
	"crypto/rand"
	"encoding/hex"
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
)

// SecurityUtils 安全工具集
type SecurityUtils struct {
	htmlPolicy *bluemonday.Policy
}

// NewSecurityUtils 创建安全工具实例
func NewSecurityUtils() *SecurityUtils {
	return &SecurityUtils{
		htmlPolicy: bluemonday.UGCPolicy(),
	}
}

// SanitizeHTML 清理HTML内容，防止XSS攻击
func (s *SecurityUtils) SanitizeHTML(input string) string {
	return s.htmlPolicy.Sanitize(input)
}

// StripHTML 完全移除HTML标签
func (s *SecurityUtils) StripHTML(input string) string {
	policy := bluemonday.StrictPolicy()
	return policy.Sanitize(input)
}

// EscapeHTML HTML转义
func (s *SecurityUtils) EscapeHTML(input string) string {
	return html.EscapeString(input)
}

// ValidateInput 综合输入验证
func (s *SecurityUtils) ValidateInput(input string, rules ValidationRules) error {
	// 长度检查
	if rules.MinLength > 0 && len(input) < rules.MinLength {
		return NewValidationError("输入长度不足")
	}
	if rules.MaxLength > 0 && len(input) > rules.MaxLength {
		return NewValidationError("输入长度超限")
	}

	// 格式检查
	if rules.Format != "" {
		switch rules.Format {
		case "email":
			if !govalidator.IsEmail(input) {
				return NewValidationError("邮箱格式无效")
			}
		case "url":
			if !govalidator.IsURL(input) {
				return NewValidationError("URL格式无效")
			}
		case "alphanumeric":
			if !govalidator.IsAlphanumeric(input) {
				return NewValidationError("只允许字母和数字")
			}
		case "alpha":
			if !govalidator.IsAlpha(input) {
				return NewValidationError("只允许字母")
			}
		case "numeric":
			if !govalidator.IsNumeric(input) {
				return NewValidationError("只允许数字")
			}
		}
	}

	// 自定义正则检查
	if rules.Pattern != "" {
		if matched, _ := regexp.MatchString(rules.Pattern, input); !matched {
			return NewValidationError("格式不符合要求")
		}
	}

	// 危险内容检查
	if s.containsDangerousContent(input) {
		return NewValidationError("包含危险内容")
	}

	return nil
}

// ValidationRules 验证规则
type ValidationRules struct {
	MinLength int    // 最小长度
	MaxLength int    // 最大长度
	Format    string // 格式类型：email, url, alphanumeric, alpha, numeric
	Pattern   string // 自定义正则表达式
	Required  bool   // 是否必填
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError 创建验证错误
func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

// containsDangerousContent 检查是否包含危险内容
func (s *SecurityUtils) containsDangerousContent(input string) bool {
	dangerousPatterns := []string{
		`<script[^>]*>.*?</script>`,
		`javascript:`,
		`vbscript:`,
		`onload\s*=`,
		`onerror\s*=`,
		`onclick\s*=`,
		`onmouseover\s*=`,
		`<iframe[^>]*>`,
		`<object[^>]*>`,
		`<embed[^>]*>`,
		`<link[^>]*>`,
		`<meta[^>]*>`,
		`expression\s*\(`,
		`url\s*\(`,
		`@import`,
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if matched, _ := regexp.MatchString(pattern, inputLower); matched {
			return true
		}
	}

	return false
}

// GenerateSecureToken 生成安全令牌
func (s *SecurityUtils) GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsValidUTF8 检查是否为有效的UTF-8编码
func (s *SecurityUtils) IsValidUTF8(input string) bool {
	// 使用标准库的utf8包检查
	return utf8.ValidString(input)
}

// ContainsOnlyPrintable 检查是否只包含可打印字符
func (s *SecurityUtils) ContainsOnlyPrintable(input string) bool {
	for _, r := range input {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// NormalizeWhitespace 规范化空白字符
func (s *SecurityUtils) NormalizeWhitespace(input string) string {
	// 替换所有连续的空白字符为单个空格
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(input, " ")
	return strings.TrimSpace(normalized)
}

// ValidateProjectName 验证项目名称
func (s *SecurityUtils) ValidateProjectName(name string) error {
	rules := ValidationRules{
		MinLength: 1,
		MaxLength: 100,
		Required:  true,
	}

	if err := s.ValidateInput(name, rules); err != nil {
		return err
	}

	// 项目名称特殊规则：不能包含特殊字符
	if matched, _ := regexp.MatchString(`[<>:"\\|?*]`, name); matched {
		return NewValidationError("项目名称不能包含特殊字符")
	}

	return nil
}

// ValidateTranslationKey 验证翻译键名
func (s *SecurityUtils) ValidateTranslationKey(key string) error {
	rules := ValidationRules{
		MinLength: 1,
		MaxLength: 200,
		Pattern:   `^[a-zA-Z0-9._-]+$`, // 只允许字母、数字、点、下划线、连字符
		Required:  true,
	}

	return s.ValidateInput(key, rules)
}

// ValidateTranslationValue 验证翻译值
func (s *SecurityUtils) ValidateTranslationValue(value string) error {
	rules := ValidationRules{
		MinLength: 0,
		MaxLength: 10000, // 10KB
		Required:  false,
	}

	if err := s.ValidateInput(value, rules); err != nil {
		return err
	}

	// 检查UTF-8有效性
	if !s.IsValidUTF8(value) {
		return NewValidationError("翻译值包含无效字符")
	}

	return nil
}

// ValidateUsername 验证用户名
func (s *SecurityUtils) ValidateUsername(username string) error {
	rules := ValidationRules{
		MinLength: 3,
		MaxLength: 50,
		Pattern:   `^[a-zA-Z0-9_-]+$`, // 只允许字母、数字、下划线、连字符
		Required:  true,
	}

	return s.ValidateInput(username, rules)
}

// ValidatePassword 验证密码强度
func (s *SecurityUtils) ValidatePassword(password string) error {
	if len(password) < 8 {
		return NewValidationError("密码长度至少8位")
	}
	if len(password) > 128 {
		return NewValidationError("密码长度不能超过128位")
	}

	// 检查密码复杂度
	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return NewValidationError("密码必须包含大写字母、小写字母和数字")
	}

	return nil
}

// CleanUserInput 清理用户输入
func (s *SecurityUtils) CleanUserInput(input string) string {
	// 1. HTML清理
	cleaned := s.SanitizeHTML(input)

	// 2. 规范化空白字符
	cleaned = s.NormalizeWhitespace(cleaned)

	// 3. 移除控制字符
	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1 // 移除控制字符
		}
		return r
	}, cleaned)

	return cleaned
}

// 全局安全工具实例
var GlobalSecurityUtils = NewSecurityUtils()
