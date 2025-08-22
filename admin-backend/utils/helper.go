package utils

import (
	"strconv"
	"strings"
)

// ParseInt 安全地解析整数字符串
// 返回解析后的整数，如果解析失败返回 defaultValue
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	// 使用标准库的 strconv.Atoi 进行解析
	if result, err := strconv.Atoi(s); err == nil {
		return result
	}

	return defaultValue
}

// ParseIntWithRange 安全地解析整数字符串并验证范围
// 返回在 [min, max] 范围内的整数，超出范围返回相应的边界值
func ParseIntWithRange(s string, defaultValue, min, max int) int {
	value := ParseInt(s, defaultValue)

	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

// IsValidInteger 检查字符串是否为有效的整数
func IsValidInteger(s string) bool {
	if s == "" {
		return false
	}

	// 处理负数
	if strings.HasPrefix(s, "-") {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// 检查每个字符是否为数字
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// ParsePositiveInt 解析正整数，如果不是正整数返回 defaultValue
func ParsePositiveInt(s string, defaultValue int) int {
	value := ParseInt(s, defaultValue)
	if value <= 0 {
		return defaultValue
	}
	return value
}

// SanitizeString 清理字符串，移除前后空白并限制长度
func SanitizeString(s string, maxLength int) string {
	// 移除前后空白
	s = strings.TrimSpace(s)

	// 限制长度
	if len(s) > maxLength {
		return s[:maxLength]
	}

	return s
}

// ContainsAny 检查字符串是否包含任何一个子字符串
func ContainsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
