package dto

// CreateLanguageRequest 创建语言请求
type CreateLanguageRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	IsDefault bool   `json:"is_default"`
}
