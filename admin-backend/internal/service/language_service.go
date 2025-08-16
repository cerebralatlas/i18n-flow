package service

import (
	"context"
	"i18n-flow/internal/domain"
	"strings"
)

// LanguageService 语言服务实现
type LanguageService struct {
	languageRepo domain.LanguageRepository
}

// NewLanguageService 创建语言服务实例
func NewLanguageService(languageRepo domain.LanguageRepository) *LanguageService {
	return &LanguageService{
		languageRepo: languageRepo,
	}
}

// Create 创建语言
func (s *LanguageService) Create(ctx context.Context, req domain.CreateLanguageRequest) (*domain.Language, error) {
	// 验证语言代码格式
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, domain.ErrInvalidLanguage
	}

	// 检查语言代码是否已存在
	existingLanguage, err := s.languageRepo.GetByCode(ctx, code)
	if err == nil && existingLanguage != nil {
		return nil, domain.ErrLanguageExists
	}

	// 如果设置为默认语言，需要先取消其他默认语言
	if req.IsDefault {
		if err := s.clearDefaultLanguage(ctx); err != nil {
			return nil, err
		}
	}

	// 创建语言
	language := &domain.Language{
		Code:      code,
		Name:      strings.TrimSpace(req.Name),
		IsDefault: req.IsDefault,
		Status:    "active",
	}

	if err := s.languageRepo.Create(ctx, language); err != nil {
		return nil, err
	}

	return language, nil
}

// GetByID 根据ID获取语言
func (s *LanguageService) GetByID(ctx context.Context, id uint) (*domain.Language, error) {
	return s.languageRepo.GetByID(ctx, id)
}

// GetAll 获取所有语言
func (s *LanguageService) GetAll(ctx context.Context) ([]*domain.Language, error) {
	return s.languageRepo.GetAll(ctx)
}

// Update 更新语言
func (s *LanguageService) Update(ctx context.Context, id uint, req domain.CreateLanguageRequest) (*domain.Language, error) {
	// 获取现有语言
	language, err := s.languageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Code != "" {
		code := strings.TrimSpace(req.Code)
		if code != language.Code {
			// 检查新代码是否已存在
			existingLanguage, err := s.languageRepo.GetByCode(ctx, code)
			if err == nil && existingLanguage != nil && existingLanguage.ID != language.ID {
				return nil, domain.ErrLanguageExists
			}
			language.Code = code
		}
	}

	if req.Name != "" {
		language.Name = strings.TrimSpace(req.Name)
	}

	// 处理默认语言设置
	if req.IsDefault && !language.IsDefault {
		// 如果要设置为默认语言，先取消其他默认语言
		if err := s.clearDefaultLanguage(ctx); err != nil {
			return nil, err
		}
		language.IsDefault = true
	} else if !req.IsDefault && language.IsDefault {
		// 不允许取消默认语言，除非有其他默认语言
		language.IsDefault = false
	}

	// 保存更新
	if err := s.languageRepo.Update(ctx, language); err != nil {
		return nil, err
	}

	return language, nil
}

// Delete 删除语言
func (s *LanguageService) Delete(ctx context.Context, id uint) error {
	// 检查语言是否存在
	language, err := s.languageRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 不允许删除默认语言
	if language.IsDefault {
		return domain.ErrInvalidInput // 或者定义专门的错误
	}

	// 删除语言
	return s.languageRepo.Delete(ctx, id)
}

// clearDefaultLanguage 清除其他语言的默认设置
func (s *LanguageService) clearDefaultLanguage(ctx context.Context) error {
	defaultLanguage, err := s.languageRepo.GetDefault(ctx)
	if err == domain.ErrLanguageNotFound {
		// 没有默认语言，无需处理
		return nil
	}
	if err != nil {
		return err
	}

	// 取消默认设置
	defaultLanguage.IsDefault = false
	return s.languageRepo.Update(ctx, defaultLanguage)
}
