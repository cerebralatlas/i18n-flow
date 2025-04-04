package service

import (
	"errors"
	"i18n-flow/model"
	"i18n-flow/model/db"

	"gorm.io/gorm"
)

// LanguageService 语言服务
type LanguageService struct{}

// LanguageRequest 语言请求参数
type LanguageRequest struct {
	Code      string `json:"code" binding:"required" example:"fr"`       // 语言代码
	Name      string `json:"name" binding:"required" example:"Français"` // 语言名称
	IsDefault bool   `json:"is_default" example:"false"`                 // 是否为默认语言
}

// LanguageResponse 语言响应
type LanguageResponse struct {
	ID        uint   `json:"id" example:"1"`                           // 语言ID
	Code      string `json:"code" example:"fr"`                        // 语言代码
	Name      string `json:"name" example:"Français"`                  // 语言名称
	IsDefault bool   `json:"is_default" example:"false"`               // 是否为默认语言
	Status    string `json:"status" example:"active"`                  // 语言状态
	CreatedAt string `json:"created_at" example:"2023-04-01 12:00:00"` // 创建时间
	UpdatedAt string `json:"updated_at" example:"2023-04-01 12:00:00"` // 更新时间
}

// GetLanguages 获取所有语言
func (s *LanguageService) GetLanguages() ([]LanguageResponse, error) {
	var languages []model.Language

	if err := db.DB.Find(&languages).Error; err != nil {
		return nil, err
	}

	// 格式化响应
	var response []LanguageResponse
	for _, lang := range languages {
		response = append(response, LanguageResponse{
			ID:        lang.ID,
			Code:      lang.Code,
			Name:      lang.Name,
			IsDefault: lang.IsDefault,
			Status:    lang.Status,
			CreatedAt: lang.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: lang.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, nil
}

// CreateLanguage 创建语言
func (s *LanguageService) CreateLanguage(req LanguageRequest) (*LanguageResponse, error) {
	// 检查语言代码是否已存在
	var existingLanguage model.Language
	if err := db.DB.Where("code = ?", req.Code).First(&existingLanguage).Error; err == nil {
		return nil, errors.New("语言代码已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 如果设置为默认语言，先将所有语言设为非默认
	if req.IsDefault {
		if err := db.DB.Model(&model.Language{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return nil, err
		}
	}

	language := model.Language{
		Code:      req.Code,
		Name:      req.Name,
		IsDefault: req.IsDefault,
		Status:    "active",
	}

	if err := db.DB.Create(&language).Error; err != nil {
		return nil, err
	}

	return &LanguageResponse{
		ID:        language.ID,
		Code:      language.Code,
		Name:      language.Name,
		IsDefault: language.IsDefault,
		Status:    language.Status,
		CreatedAt: language.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: language.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateLanguage 更新语言
func (s *LanguageService) UpdateLanguage(id uint, req LanguageRequest) (*LanguageResponse, error) {
	var language model.Language

	if err := db.DB.First(&language, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("语言不存在")
		}
		return nil, err
	}

	// 如果语言代码变更，检查是否与其他语言冲突
	if req.Code != language.Code {
		var existingLanguage model.Language
		if err := db.DB.Where("code = ? AND id != ?", req.Code, id).First(&existingLanguage).Error; err == nil {
			return nil, errors.New("语言代码已存在")
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	// 如果设置为默认语言，先将所有语言设为非默认
	if req.IsDefault && !language.IsDefault {
		if err := db.DB.Model(&model.Language{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return nil, err
		}
	}

	// 更新语言信息
	language.Code = req.Code
	language.Name = req.Name
	language.IsDefault = req.IsDefault

	if err := db.DB.Save(&language).Error; err != nil {
		return nil, err
	}

	return &LanguageResponse{
		ID:        language.ID,
		Code:      language.Code,
		Name:      language.Name,
		IsDefault: language.IsDefault,
		Status:    language.Status,
		CreatedAt: language.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: language.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// DeleteLanguage 删除语言
func (s *LanguageService) DeleteLanguage(id uint) error {
	var language model.Language

	// 获取语言信息
	if err := db.DB.First(&language, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("语言不存在")
		}
		return err
	}

	// 检查是否为默认语言
	if language.IsDefault {
		return errors.New("不能删除默认语言")
	}

	// 检查是否有关联的翻译
	var count int64
	if err := db.DB.Model(&model.Translation{}).Where("language_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("该语言有关联的翻译，无法删除")
	}

	// 删除语言
	if err := db.DB.Delete(&language).Error; err != nil {
		return err
	}

	return nil
}
