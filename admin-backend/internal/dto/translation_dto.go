package dto

// CreateTranslationRequest 创建翻译请求
type CreateTranslationRequest struct {
	ProjectID  uint64 `json:"project_id" binding:"required"`
	KeyName    string `json:"key_name" binding:"required"`
	Context    string `json:"context"`
	LanguageID uint64 `json:"language_id" binding:"required"`
	Value      string `json:"value" binding:"required"`
}

// BatchTranslationRequest 批量翻译请求（前端格式）
type BatchTranslationRequest struct {
	ProjectID    uint64            `json:"project_id" binding:"required"`
	KeyName      string            `json:"key_name" binding:"required"`
	Context      string            `json:"context"`
	Translations map[string]string `json:"translations" binding:"required"`
}
