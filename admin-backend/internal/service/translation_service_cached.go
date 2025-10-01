package service

import (
	"context"
	"encoding/json"
	"fmt"
	"i18n-flow/internal/domain"
	"strconv"
	"sync"
	"time"
)

// CachedTranslationService 带缓存的翻译服务实现
type CachedTranslationService struct {
	translationService *TranslationService
	cacheService       domain.CacheService
	// 用于防止缓存击穿的互斥锁
	cacheMutexes map[string]*sync.Mutex
	mutexLock    sync.RWMutex
}

// NewCachedTranslationService 创建带缓存的翻译服务实例
func NewCachedTranslationService(
	translationService *TranslationService,
	cacheService domain.CacheService,
) *CachedTranslationService {
	return &CachedTranslationService{
		translationService: translationService,
		cacheService:       cacheService,
		cacheMutexes:       make(map[string]*sync.Mutex),
	}
}

// getMutex 获取指定键的互斥锁，用于防止缓存击穿
func (s *CachedTranslationService) getMutex(key string) *sync.Mutex {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()

	if mutex, exists := s.cacheMutexes[key]; exists {
		return mutex
	}

	mutex := &sync.Mutex{}
	s.cacheMutexes[key] = mutex
	return mutex
}

// removeMutex 移除指定键的互斥锁
func (s *CachedTranslationService) removeMutex(key string) {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()

	delete(s.cacheMutexes, key)
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

// UpsertBatch 批量创建或更新翻译（更新缓存）
func (s *CachedTranslationService) UpsertBatch(ctx context.Context, requests []domain.CreateTranslationRequest) error {
	err := s.translationService.UpsertBatch(ctx, requests)
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

// GetByID 根据ID获取翻译
func (s *CachedTranslationService) GetByID(ctx context.Context, id uint) (*domain.Translation, error) {
	// 这个方法不缓存，因为单个翻译查询不频繁
	return s.translationService.GetByID(ctx, id)
}

// TranslationCacheResult 定义翻译缓存结果结构体
type TranslationCacheResult struct {
	Translations []*domain.Translation `json:"translations"`
	Total        int64                 `json:"total"`
}

// GetByProjectID 根据项目ID获取翻译（使用缓存）
func (s *CachedTranslationService) GetByProjectID(ctx context.Context, projectID uint, limit, offset int) ([]*domain.Translation, int64, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("%s:%d:%d", s.cacheService.GetTranslationKey(projectID), limit, offset)

	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	var cachedResult TranslationCacheResult
	err := s.cacheService.GetJSONWithEmptyCheck(ctx, cacheKey, &cachedResult)
	if err == nil {
		//fmt.Printf("翻译列表缓存命中 [project=%d, total=%d]\n", projectID, cachedResult.Total)
		return cachedResult.Translations, cachedResult.Total, nil
	}

	// 缓存未命中，从数据库获取
	translations, total, err := s.translationService.GetByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// 更新缓存，添加随机过期时间防止雪崩
	cachedResult = TranslationCacheResult{
		Translations: translations,
		Total:        total,
	}

	expiration := s.cacheService.AddRandomExpiration(domain.DefaultExpiration)
	if err := s.cacheService.SetJSONWithEmptyCache(ctx, cacheKey, cachedResult, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果，记录日志用于调试
		//fmt.Printf("翻译列表缓存更新失败 [project=%d, limit=%d, offset=%d]: %v\n", projectID, limit, offset, err)
	} /*else {
		fmt.Printf("翻译列表缓存更新成功 [project=%d, total=%d]\n", projectID, total)
	}*/

	return translations, total, nil
}

// MatrixCacheResult 定义缓存结果结构体
type MatrixCacheResult struct {
	Matrix map[string]map[string]string `json:"matrix"`
	Total  int64                        `json:"total"`
}

// GetMatrix 获取翻译矩阵（使用缓存）
func (s *CachedTranslationService) GetMatrix(ctx context.Context, projectID uint, limit, offset int, keyword string) (map[string]map[string]string, int64, error) {
	// 优化缓存键生成，区分搜索和非搜索查询
	var cacheKey string
	if keyword != "" {
		// 搜索查询使用较短的缓存时间
		cacheKey = fmt.Sprintf("%s:search:%s:%d:%d", s.cacheService.GetTranslationMatrixKey(projectID, ""), s.hashKeyword(keyword), limit, offset)
	} else {
		// 非搜索查询使用较长的缓存时间
		cacheKey = fmt.Sprintf("%s:all:%d:%d", s.cacheService.GetTranslationMatrixKey(projectID, ""), limit, offset)
	}

	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	var cachedResult MatrixCacheResult
	err := s.cacheService.GetJSONWithEmptyCheck(ctx, cacheKey, &cachedResult)
	if err == nil {
		//fmt.Printf("翻译矩阵缓存命中 [project=%d, keyword=%s, total=%d]\n", projectID, keyword, cachedResult.Total)
		return cachedResult.Matrix, cachedResult.Total, nil
	}

	// 缓存未命中，从数据库获取
	matrix, total, err := s.translationService.GetMatrix(ctx, projectID, limit, offset, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 更新缓存，添加随机过期时间防止雪崩
	cachedResult = MatrixCacheResult{
		Matrix: matrix,
		Total:  total,
	}

	// 根据查询类型设置不同的缓存时间
	var expiration time.Duration
	if keyword != "" {
		// 搜索查询缓存较短时间
		expiration = s.cacheService.AddRandomExpiration(5 * time.Minute)
	} else {
		// 非搜索查询缓存较长时间
		expiration = s.cacheService.AddRandomExpiration(domain.DefaultExpiration)
	}

	if err := s.cacheService.SetJSONWithEmptyCache(ctx, cacheKey, cachedResult, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果，记录日志用于调试
		//fmt.Printf("翻译矩阵缓存更新失败 [project=%d, keyword=%s, limit=%d, offset=%d]: %v\n", projectID, keyword, limit, offset, err)
	} /*else {
		fmt.Printf("翻译矩阵缓存更新成功 [project=%d, keyword=%s, total=%d, keys=%d]\n", projectID, keyword, total, len(matrix))
	}*/

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
	// 使用管道操作提高性能
	// 清除翻译列表缓存
	s.cacheService.DeleteByPattern(ctx, s.cacheService.GetTranslationKey(projectID)+"*")

	// 清除翻译矩阵缓存
	s.cacheService.DeleteByPattern(ctx, s.cacheService.GetTranslationMatrixKey(projectID, "")+"*")

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())
}

// invalidateLanguageCache 清除语言相关的缓存（当语言被修改时调用）
func (s *CachedTranslationService) invalidateLanguageCache(ctx context.Context) {
	// 清除所有项目的翻译矩阵缓存，因为语言变更可能影响所有项目
	s.cacheService.DeleteByPattern(ctx, domain.TranslationMatrixPrefix+"*")

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())
}

// invalidateSpecificTranslationCache 清除特定翻译键的缓存
func (s *CachedTranslationService) invalidateSpecificTranslationCache(ctx context.Context, projectID uint, keyName string) {
	// 清除包含特定键名的翻译矩阵缓存
	pattern := fmt.Sprintf("%s%d:*%s*", domain.TranslationMatrixPrefix, projectID, keyName)
	s.cacheService.DeleteByPattern(ctx, pattern)

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())
}

// hashKeyword 对关键词进行简单哈希，避免缓存键过长
func (s *CachedTranslationService) hashKeyword(keyword string) string {
	// 简单的哈希函数，生产环境可以使用更复杂的哈希
	hash := 0
	for _, char := range keyword {
		hash = 31*hash + int(char)
	}
	return strconv.Itoa(hash)
}
