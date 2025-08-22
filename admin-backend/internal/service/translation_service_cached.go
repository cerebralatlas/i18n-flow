package service

import (
	"context"
	"encoding/json"
	"fmt"
	"i18n-flow/internal/domain"
)

// CachedTranslationService 带缓存的翻译服务实现
type CachedTranslationService struct {
	translationService *TranslationService
	cacheService       domain.CacheService
}

// NewCachedTranslationService 创建带缓存的翻译服务实例
func NewCachedTranslationService(
	translationService *TranslationService,
	cacheService domain.CacheService,
) *CachedTranslationService {
	return &CachedTranslationService{
		translationService: translationService,
		cacheService:       cacheService,
	}
}

// Create 创建翻译（更新缓存）
func (s *CachedTranslationService) Create(ctx context.Context, req domain.CreateTranslationRequest) (*domain.Translation, error) {
	translation, err := s.translationService.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// 清除相关缓存
	s.invalidateProjectCache(ctx, req.ProjectID)

	return translation, nil
}

// CreateBatch 批量创建翻译（更新缓存）
func (s *CachedTranslationService) CreateBatch(ctx context.Context, requests []domain.CreateTranslationRequest) error {
	err := s.translationService.CreateBatch(ctx, requests)
	if err != nil {
		return err
	}

	// 清除相关缓存
	projectIDs := make(map[uint]bool)
	for _, req := range requests {
		projectIDs[req.ProjectID] = true
	}

	for projectID := range projectIDs {
		s.invalidateProjectCache(ctx, projectID)
	}

	return nil
}

// CreateBatchFromRequest 从批量翻译请求创建翻译（更新缓存）
func (s *CachedTranslationService) CreateBatchFromRequest(ctx context.Context, req domain.BatchTranslationRequest) error {
	err := s.translationService.CreateBatchFromRequest(ctx, req)
	if err != nil {
		return err
	}

	// 清除相关缓存
	s.invalidateProjectCache(ctx, req.ProjectID)

	return nil
}

// GetByID 根据ID获取翻译
func (s *CachedTranslationService) GetByID(ctx context.Context, id uint) (*domain.Translation, error) {
	// 这个方法不缓存，因为单个翻译查询不频繁
	return s.translationService.GetByID(ctx, id)
}

// GetByProjectID 根据项目ID获取翻译（使用缓存）
func (s *CachedTranslationService) GetByProjectID(ctx context.Context, projectID uint, limit, offset int) ([]*domain.Translation, int64, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("%s:%d:%d", s.cacheService.GetTranslationKey(projectID), limit, offset)

	// 尝试从缓存获取
	var cachedResult struct {
		Translations []*domain.Translation `json:"translations"`
		Total        int64                 `json:"total"`
	}

	err := s.cacheService.GetJSON(ctx, cacheKey, &cachedResult)
	if err == nil {
		return cachedResult.Translations, cachedResult.Total, nil
	}

	// 缓存未命中，从数据库获取
	translations, total, err := s.translationService.GetByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// 更新缓存
	cachedResult = struct {
		Translations []*domain.Translation `json:"translations"`
		Total        int64                 `json:"total"`
	}{
		Translations: translations,
		Total:        total,
	}

	if err := s.cacheService.SetJSON(ctx, cacheKey, cachedResult, domain.DefaultExpiration); err != nil {
		// 缓存更新失败，但不影响返回结果
		fmt.Printf("缓存更新失败: %v\n", err)
	}

	return translations, total, nil
}

// GetMatrix 获取翻译矩阵（使用缓存）
func (s *CachedTranslationService) GetMatrix(ctx context.Context, projectID uint, limit, offset int, keyword string) (map[string]map[string]string, int64, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("%s:%d:%d", s.cacheService.GetTranslationMatrixKey(projectID, keyword), limit, offset)

	// 尝试从缓存获取
	var cachedResult struct {
		Matrix map[string]map[string]string `json:"matrix"`
		Total  int64                        `json:"total"`
	}

	err := s.cacheService.GetJSON(ctx, cacheKey, &cachedResult)
	if err == nil {
		return cachedResult.Matrix, cachedResult.Total, nil
	}

	// 缓存未命中，从数据库获取
	matrix, total, err := s.translationService.GetMatrix(ctx, projectID, limit, offset, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 更新缓存
	cachedResult = struct {
		Matrix map[string]map[string]string `json:"matrix"`
		Total  int64                        `json:"total"`
	}{
		Matrix: matrix,
		Total:  total,
	}

	if err := s.cacheService.SetJSON(ctx, cacheKey, cachedResult, domain.DefaultExpiration); err != nil {
		// 缓存更新失败，但不影响返回结果
		fmt.Printf("缓存更新失败: %v\n", err)
	}

	return matrix, total, nil
}

// Update 更新翻译（更新缓存）
func (s *CachedTranslationService) Update(ctx context.Context, id uint, req domain.CreateTranslationRequest) (*domain.Translation, error) {
	// 先获取原始翻译，用于后续清除缓存
	oldTranslation, err := s.translationService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	translation, err := s.translationService.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// 清除相关缓存
	s.invalidateProjectCache(ctx, oldTranslation.ProjectID)
	if req.ProjectID != 0 && req.ProjectID != oldTranslation.ProjectID {
		s.invalidateProjectCache(ctx, req.ProjectID)
	}

	return translation, nil
}

// Delete 删除翻译（更新缓存）
func (s *CachedTranslationService) Delete(ctx context.Context, id uint) error {
	// 先获取翻译，用于后续清除缓存
	translation, err := s.translationService.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.translationService.Delete(ctx, id)
	if err != nil {
		return err
	}

	// 清除相关缓存
	s.invalidateProjectCache(ctx, translation.ProjectID)

	return nil
}

// DeleteBatch 批量删除翻译（更新缓存）
func (s *CachedTranslationService) DeleteBatch(ctx context.Context, ids []uint) error {
	// 这里需要先查询所有翻译，获取相关的项目ID
	projectIDs := make(map[uint]bool)
	for _, id := range ids {
		translation, err := s.translationService.GetByID(ctx, id)
		if err == nil {
			projectIDs[translation.ProjectID] = true
		}
	}

	err := s.translationService.DeleteBatch(ctx, ids)
	if err != nil {
		return err
	}

	// 清除相关缓存
	for projectID := range projectIDs {
		s.invalidateProjectCache(ctx, projectID)
	}

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())

	return nil
}

// Export 导出翻译
func (s *CachedTranslationService) Export(ctx context.Context, projectID uint, format string) ([]byte, error) {
	// 使用缓存的矩阵数据
	matrix, _, err := s.GetMatrix(ctx, projectID, -1, 0, "")
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

// Import 导入翻译（更新缓存）
func (s *CachedTranslationService) Import(ctx context.Context, projectID uint, data []byte, format string) error {
	err := s.translationService.Import(ctx, projectID, data, format)
	if err != nil {
		return err
	}

	// 清除相关缓存
	s.invalidateProjectCache(ctx, projectID)

	return nil
}

// invalidateProjectCache 清除项目相关的所有缓存
func (s *CachedTranslationService) invalidateProjectCache(ctx context.Context, projectID uint) {
	// 清除翻译列表缓存
	s.cacheService.DeleteByPattern(ctx, s.cacheService.GetTranslationKey(projectID)+"*")

	// 清除翻译矩阵缓存
	s.cacheService.DeleteByPattern(ctx, s.cacheService.GetTranslationMatrixKey(projectID, "")+"*")

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())
}
