package repository

import (
	"context"
	"errors"
	"i18n-flow/internal/domain"

	"gorm.io/gorm"
)

// TranslationRepository 翻译仓储实现
type TranslationRepository struct {
	db *gorm.DB
}

// NewTranslationRepository 创建翻译仓储实例
func NewTranslationRepository(db *gorm.DB) *TranslationRepository {
	return &TranslationRepository{db: db}
}

// GetByID 根据ID获取翻译
func (r *TranslationRepository) GetByID(ctx context.Context, id uint) (*domain.Translation, error) {
	var translation domain.Translation
	if err := r.db.WithContext(ctx).Preload("Project").Preload("Language").First(&translation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTranslationNotFound
		}
		return nil, err
	}
	return &translation, nil
}

// GetByProjectID 根据项目ID获取翻译（分页）
func (r *TranslationRepository) GetByProjectID(ctx context.Context, projectID uint, limit, offset int) ([]*domain.Translation, int64, error) {
	var translations []*domain.Translation
	var total int64

	query := r.db.WithContext(ctx).Where("project_id = ?", projectID)

	// 计算总数
	if err := query.Model(&domain.Translation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Preload("Language").Limit(limit).Offset(offset).Find(&translations).Error; err != nil {
		return nil, 0, err
	}

	return translations, total, nil
}

// GetByProjectAndLanguage 根据项目和语言获取翻译
func (r *TranslationRepository) GetByProjectAndLanguage(ctx context.Context, projectID, languageID uint) ([]*domain.Translation, error) {
	var translations []*domain.Translation
	if err := r.db.WithContext(ctx).Where("project_id = ? AND language_id = ?", projectID, languageID).Find(&translations).Error; err != nil {
		return nil, err
	}
	return translations, nil
}

// GetMatrix 获取翻译矩阵（key-language映射）
func (r *TranslationRepository) GetMatrix(ctx context.Context, projectID uint) (map[string]map[string]string, error) {
	var translations []*domain.Translation
	if err := r.db.WithContext(ctx).Preload("Language").Where("project_id = ?", projectID).Find(&translations).Error; err != nil {
		return nil, err
	}

	matrix := make(map[string]map[string]string)
	for _, translation := range translations {
		if matrix[translation.KeyName] == nil {
			matrix[translation.KeyName] = make(map[string]string)
		}
		matrix[translation.KeyName][translation.Language.Code] = translation.Value
	}

	return matrix, nil
}

// Create 创建翻译
func (r *TranslationRepository) Create(ctx context.Context, translation *domain.Translation) error {
	return r.db.WithContext(ctx).Create(translation).Error
}

// CreateBatch 批量创建翻译
func (r *TranslationRepository) CreateBatch(ctx context.Context, translations []*domain.Translation) error {
	if len(translations) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(translations, 100).Error
}

// Update 更新翻译
func (r *TranslationRepository) Update(ctx context.Context, translation *domain.Translation) error {
	return r.db.WithContext(ctx).Save(translation).Error
}

// Delete 删除翻译
func (r *TranslationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Translation{}, id).Error
}

// DeleteBatch 批量删除翻译
func (r *TranslationRepository) DeleteBatch(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&domain.Translation{}, ids).Error
}
