package service

import (
	"context"
	"encoding/json"
	"fmt"
	"i18n-flow/internal/domain"
	"strings"
)

// TranslationService 翻译服务实现
type TranslationService struct {
	translationRepo domain.TranslationRepository
	projectRepo     domain.ProjectRepository
	languageRepo    domain.LanguageRepository
}

// NewTranslationService 创建翻译服务实例
func NewTranslationService(
	translationRepo domain.TranslationRepository,
	projectRepo domain.ProjectRepository,
	languageRepo domain.LanguageRepository,
) *TranslationService {
	return &TranslationService{
		translationRepo: translationRepo,
		projectRepo:     projectRepo,
		languageRepo:    languageRepo,
	}
}

// Create 创建翻译
func (s *TranslationService) Create(ctx context.Context, req domain.CreateTranslationRequest) (*domain.Translation, error) {
	// 验证项目是否存在
	_, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, domain.ErrProjectNotFound
	}

	// 验证语言是否存在
	_, err = s.languageRepo.GetByID(ctx, req.LanguageID)
	if err != nil {
		return nil, domain.ErrLanguageNotFound
	}

	// 检查翻译是否已存在
	keyName := strings.TrimSpace(req.KeyName)
	existing, err := s.translationRepo.GetByProjectKeyLanguage(ctx, req.ProjectID, keyName, req.LanguageID)
	if err == nil && existing != nil {
		return nil, domain.NewAppErrorWithDetails(
			domain.ErrorTypeConflict,
			"TRANSLATION_EXISTS",
			"该项目中已存在相同键名和语言的翻译",
			fmt.Sprintf("项目ID: %d, 键名: %s, 语言ID: %d", req.ProjectID, keyName, req.LanguageID),
		)
	}

	// 创建翻译
	translation := &domain.Translation{
		ProjectID:  req.ProjectID,
		KeyName:    keyName,
		Context:    strings.TrimSpace(req.Context),
		LanguageID: req.LanguageID,
		Value:      strings.TrimSpace(req.Value),
		Status:     "active",
	}

	if err := s.translationRepo.Create(ctx, translation); err != nil {
		// 检查是否是唯一约束冲突错误
		if isDuplicateKeyError(err) {
			return nil, domain.NewAppErrorWithDetails(
				domain.ErrorTypeConflict,
				"TRANSLATION_EXISTS",
				"该项目中已存在相同键名和语言的翻译",
				fmt.Sprintf("项目ID: %d, 键名: %s, 语言ID: %d", req.ProjectID, keyName, req.LanguageID),
			)
		}
		return nil, err
	}

	return translation, nil
}

// CreateBatch 批量创建翻译
func (s *TranslationService) CreateBatch(ctx context.Context, requests []domain.CreateTranslationRequest) error {
	if len(requests) == 0 {
		return nil
	}

	// 收集所有请求中的项目和语言ID
	projectIDSet := make(map[uint]bool)
	languageIDSet := make(map[uint]bool)

	for _, req := range requests {
		projectIDSet[req.ProjectID] = true
		languageIDSet[req.LanguageID] = true
	}

	// 转换为切片
	projectIDs := make([]uint, 0, len(projectIDSet))
	for id := range projectIDSet {
		projectIDs = append(projectIDs, id)
	}
	languageIDs := make([]uint, 0, len(languageIDSet))
	for id := range languageIDSet {
		languageIDs = append(languageIDs, id)
	}

	// 批量验证项目 (修复 N+1 查询)
	projects, err := s.projectRepo.GetByIDs(ctx, projectIDs)
	if err != nil {
		return err
	}
	if len(projects) != len(projectIDs) {
		return domain.ErrProjectNotFound
	}

	// 批量验证语言 (修复 N+1 查询)
	languages, err := s.languageRepo.GetByIDs(ctx, languageIDs)
	if err != nil {
		return err
	}
	if len(languages) != len(languageIDs) {
		return domain.ErrLanguageNotFound
	}

	// 检查重复翻译并转换为domain对象
	translations := make([]*domain.Translation, 0, len(requests))
	duplicates := make([]string, 0)

	for _, req := range requests {
		keyName := strings.TrimSpace(req.KeyName)

		// 检查是否已存在
		existing, err := s.translationRepo.GetByProjectKeyLanguage(ctx, req.ProjectID, keyName, req.LanguageID)
		if err != nil {
			return err
		}

		if existing != nil {
			duplicates = append(duplicates, fmt.Sprintf("项目ID:%d, 键名:%s, 语言ID:%d", req.ProjectID, keyName, req.LanguageID))
			continue
		}

		translations = append(translations, &domain.Translation{
			ProjectID:  req.ProjectID,
			KeyName:    keyName,
			Context:    strings.TrimSpace(req.Context),
			LanguageID: req.LanguageID,
			Value:      strings.TrimSpace(req.Value),
			Status:     "active",
		})
	}

	// 如果有重复项，返回错误
	if len(duplicates) > 0 {
		return domain.NewAppErrorWithDetails(
			domain.ErrorTypeConflict,
			"TRANSLATION_EXISTS",
			"批量创建中存在重复的翻译",
			fmt.Sprintf("重复项: %s", strings.Join(duplicates, "; ")),
		)
	}

	// 如果没有有效的翻译需要创建
	if len(translations) == 0 {
		return nil
	}

	return s.translationRepo.CreateBatch(ctx, translations)
}

// UpsertBatch 批量创建或更新翻译
// 如果翻译已存在（基于 project_id + key_name + language_id），则更新
// 如果不存在，则创建
func (s *TranslationService) UpsertBatch(ctx context.Context, requests []domain.CreateTranslationRequest) error {
	if len(requests) == 0 {
		return nil
	}

	// 收集所有请求中的项目和语言ID
	projectIDSet := make(map[uint]bool)
	languageIDSet := make(map[uint]bool)

	for _, req := range requests {
		projectIDSet[req.ProjectID] = true
		languageIDSet[req.LanguageID] = true
	}

	// 转换为切片
	projectIDs := make([]uint, 0, len(projectIDSet))
	for id := range projectIDSet {
		projectIDs = append(projectIDs, id)
	}
	languageIDs := make([]uint, 0, len(languageIDSet))
	for id := range languageIDSet {
		languageIDs = append(languageIDs, id)
	}

	// 批量验证项目 (修复 N+1 查询)
	projects, err := s.projectRepo.GetByIDs(ctx, projectIDs)
	if err != nil {
		return err
	}
	if len(projects) != len(projectIDs) {
		return domain.ErrProjectNotFound
	}

	// 批量验证语言 (修复 N+1 查询)
	languages, err := s.languageRepo.GetByIDs(ctx, languageIDs)
	if err != nil {
		return err
	}
	if len(languages) != len(languageIDs) {
		return domain.ErrLanguageNotFound
	}

	// 转换为 domain 对象
	translations := make([]*domain.Translation, 0, len(requests))
	for _, req := range requests {
		translations = append(translations, &domain.Translation{
			ProjectID:  req.ProjectID,
			KeyName:    strings.TrimSpace(req.KeyName),
			Context:    strings.TrimSpace(req.Context),
			LanguageID: req.LanguageID,
			Value:      strings.TrimSpace(req.Value),
			Status:     "active",
		})
	}

	// 使用 UpsertBatch 而不是 CreateBatch
	return s.translationRepo.UpsertBatch(ctx, translations)
}

// CreateBatchFromRequest 从批量翻译请求创建或更新翻译
// 现在使用 UpsertBatch，支持创建和更新操作
func (s *TranslationService) CreateBatchFromRequest(ctx context.Context, req domain.BatchTranslationRequest) error {
	// 获取所有语言
	languages, err := s.languageRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// 创建语言代码到ID的映射
	languageCodeToID := make(map[string]uint)
	for _, lang := range languages {
		languageCodeToID[lang.Code] = lang.ID
	}

	// 转换为标准翻译请求
	var requests []domain.CreateTranslationRequest
	for langCode, value := range req.Translations {
		// 跳过空值
		if value == "" {
			continue
		}

		if langID, exists := languageCodeToID[langCode]; exists {
			requests = append(requests, domain.CreateTranslationRequest{
				ProjectID:  req.ProjectID,
				KeyName:    req.KeyName,
				Context:    req.Context,
				LanguageID: langID,
				Value:      value,
			})
		}
	}

	if len(requests) == 0 {
		return fmt.Errorf("no valid translations to create")
	}

	// 使用 UpsertBatch 而不是 CreateBatch，支持创建和更新
	return s.UpsertBatch(ctx, requests)
}

// GetByID 根据ID获取翻译
func (s *TranslationService) GetByID(ctx context.Context, id uint) (*domain.Translation, error) {
	return s.translationRepo.GetByID(ctx, id)
}

// GetByProjectID 根据项目ID获取翻译
func (s *TranslationService) GetByProjectID(ctx context.Context, projectID uint, limit, offset int) ([]*domain.Translation, int64, error) {
	// 验证项目是否存在
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, 0, domain.ErrProjectNotFound
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.translationRepo.GetByProjectID(ctx, projectID, limit, offset)
}

// GetMatrix 获取翻译矩阵
func (s *TranslationService) GetMatrix(ctx context.Context, projectID uint, limit, offset int, keyword string) (map[string]map[string]string, int64, error) {
	// 验证项目是否存在
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, 0, domain.ErrProjectNotFound
	}

	return s.translationRepo.GetMatrix(ctx, projectID, limit, offset, keyword)
}

// Update 更新翻译
func (s *TranslationService) Update(ctx context.Context, id uint, req domain.CreateTranslationRequest) (*domain.Translation, error) {
	// 获取现有翻译
	translation, err := s.translationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 如果项目ID改变，验证新项目
	if req.ProjectID != 0 && req.ProjectID != translation.ProjectID {
		_, err := s.projectRepo.GetByID(ctx, req.ProjectID)
		if err != nil {
			return nil, domain.ErrProjectNotFound
		}
		translation.ProjectID = req.ProjectID
	}

	// 如果语言ID改变，验证新语言
	if req.LanguageID != 0 && req.LanguageID != translation.LanguageID {
		_, err := s.languageRepo.GetByID(ctx, req.LanguageID)
		if err != nil {
			return nil, domain.ErrLanguageNotFound
		}
		translation.LanguageID = req.LanguageID
	}

	// 更新其他字段
	if req.KeyName != "" {
		translation.KeyName = strings.TrimSpace(req.KeyName)
	}

	if req.Context != "" {
		translation.Context = strings.TrimSpace(req.Context)
	}

	if req.Value != "" {
		translation.Value = strings.TrimSpace(req.Value)
	}

	// 保存更新
	if err := s.translationRepo.Update(ctx, translation); err != nil {
		return nil, err
	}

	return translation, nil
}

// Delete 删除翻译
func (s *TranslationService) Delete(ctx context.Context, id uint) error {
	// 检查翻译是否存在
	_, err := s.translationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.translationRepo.Delete(ctx, id)
}

// DeleteBatch 批量删除翻译
func (s *TranslationService) DeleteBatch(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return s.translationRepo.DeleteBatch(ctx, ids)
}

// Export 导出翻译
func (s *TranslationService) Export(ctx context.Context, projectID uint, format string) ([]byte, error) {
	// 验证项目是否存在
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, domain.ErrProjectNotFound
	}

	// 获取翻译矩阵（导出所有数据，不分页）
	matrix, _, err := s.translationRepo.GetMatrix(ctx, projectID, -1, 0, "")
	if err != nil {
		return nil, err
	}

	switch format {
	case "json":
		return json.MarshalIndent(matrix, "", "  ")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// Import 导入翻译
func (s *TranslationService) Import(ctx context.Context, projectID uint, data []byte, format string) error {
	// 验证项目是否存在
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return domain.ErrProjectNotFound
	}

	switch format {
	case "json":
		return s.importFromJSON(ctx, projectID, data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// importFromJSON 从JSON导入翻译
func (s *TranslationService) importFromJSON(ctx context.Context, projectID uint, data []byte) error {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// 获取所有语言
	languages, err := s.languageRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// 创建语言代码到ID的映射
	languageCodeToID := make(map[string]uint)
	for _, lang := range languages {
		languageCodeToID[lang.Code] = lang.ID
	}

	// 转换为翻译请求
	var requests []domain.CreateTranslationRequest

	// 检测数据格式并转换
	matrix := s.normalizeImportData(rawData)

	for key, translations := range matrix {
		for langCode, value := range translations {
			if langID, exists := languageCodeToID[langCode]; exists {
				requests = append(requests, domain.CreateTranslationRequest{
					ProjectID:  projectID,
					KeyName:    key,
					LanguageID: langID,
					Value:      value,
				})
			}
		}
	}

	if len(requests) == 0 {
		return fmt.Errorf("no valid translations found in import data")
	}

	return s.CreateBatch(ctx, requests)
}

// normalizeImportData 标准化导入数据格式
// 支持两种格式：
// 1. key -> {language: value} (标准格式)
// 2. language -> {key: value} (前端格式)
func (s *TranslationService) normalizeImportData(rawData map[string]interface{}) map[string]map[string]string {
	matrix := make(map[string]map[string]string)

	// 检测数据格式
	if s.isLanguageToKeyFormat(rawData) {
		// 前端格式: language -> {key: value}
		for langCode, keysInterface := range rawData {
			if keys, ok := keysInterface.(map[string]interface{}); ok {
				for key, valueInterface := range keys {
					if value, ok := valueInterface.(string); ok {
						if matrix[key] == nil {
							matrix[key] = make(map[string]string)
						}
						matrix[key][langCode] = value
					}
				}
			}
		}
	} else {
		// 标准格式: key -> {language: value}
		for key, languagesInterface := range rawData {
			if languages, ok := languagesInterface.(map[string]interface{}); ok {
				matrix[key] = make(map[string]string)
				for langCode, valueInterface := range languages {
					if value, ok := valueInterface.(string); ok {
						matrix[key][langCode] = value
					}
				}
			}
		}
	}

	return matrix
}

// isLanguageToKeyFormat 检测是否为 language -> {key: value} 格式
func (s *TranslationService) isLanguageToKeyFormat(rawData map[string]interface{}) bool {
	// 检查第一层的键是否看起来像语言代码
	for key := range rawData {
		// 如果键是短的字符串（1-5个字符），可能是语言代码
		if len(key) <= 5 && isLikelyLanguageCode(key) {
			return true
		}
		// 如果键包含点号，更可能是翻译键而不是语言代码
		if strings.Contains(key, ".") {
			return false
		}
	}
	return false
}

// isLikelyLanguageCode 判断字符串是否像语言代码
func isLikelyLanguageCode(code string) bool {
	// 常见的语言代码模式
	commonLanguageCodes := []string{
		"en", "zh", "ja", "ko", "fr", "de", "es", "pt", "ru", "ar", "hi", "th", "vi", "id", "ms", "tr", "it", "pl", "nl", "sv", "da", "no", "fi",
		"zh_CN", "zh_TW", "en_US", "en_GB", "pt_BR", "es_ES", "fr_FR", "de_DE",
	}

	for _, lang := range commonLanguageCodes {
		if code == lang {
			return true
		}
	}

	// 简单的启发式规则：长度为2-5的字符串，只包含字母、数字和连字符
	if len(code) >= 2 && len(code) <= 5 {
		for _, c := range code {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
				return false
			}
		}
		return true
	}

	return false
}

// isDuplicateKeyError 检查是否是重复键错误
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	// MySQL重复键错误模式
	return strings.Contains(errStr, "duplicate entry") ||
		strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "idx_translation_unique")
}
