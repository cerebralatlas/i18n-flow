package service

import (
	"context"
	"fmt"
	"i18n-flow/internal/domain"
	"i18n-flow/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService 缓存服务实现
type CacheService struct {
	redisClient *repository.RedisClient
}

// NewCacheService 创建缓存服务实例
func NewCacheService(redisClient *repository.RedisClient) *CacheService {
	return &CacheService{
		redisClient: redisClient,
	}
}

// Set 设置缓存
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.redisClient.Set(ctx, key, value, expiration)
}

// Get 获取缓存
func (s *CacheService) Get(ctx context.Context, key string) (string, error) {
	val, err := s.redisClient.Get(ctx, key)
	if err == redis.Nil {
		return "", domain.ErrCacheMiss
	}
	return val, err
}

// Delete 删除缓存
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.redisClient.Delete(ctx, key)
}

// DeleteByPattern 根据模式删除缓存
func (s *CacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	return s.redisClient.DeleteByPattern(ctx, pattern)
}

// Exists 检查缓存是否存在
func (s *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	return s.redisClient.Exists(ctx, key)
}

// SetJSON 设置JSON缓存
func (s *CacheService) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.redisClient.SetJSON(ctx, key, value, expiration)
}

// GetJSON 获取JSON缓存
func (s *CacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	err := s.redisClient.GetJSON(ctx, key, dest)
	if err == redis.Nil {
		return domain.ErrCacheMiss
	}
	return err
}

// HSet 设置哈希表字段
func (s *CacheService) HSet(ctx context.Context, key, field string, value interface{}) error {
	return s.redisClient.HSet(ctx, key, field, value)
}

// HGet 获取哈希表字段
func (s *CacheService) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := s.redisClient.HGet(ctx, key, field)
	if err == redis.Nil {
		return "", domain.ErrCacheMiss
	}
	return val, err
}

// HGetAll 获取哈希表所有字段
func (s *CacheService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	val, err := s.redisClient.HGetAll(ctx, key)
	if err == redis.Nil || len(val) == 0 {
		return nil, domain.ErrCacheMiss
	}
	return val, err
}

// HDel 删除哈希表字段
func (s *CacheService) HDel(ctx context.Context, key string, fields ...string) error {
	return s.redisClient.HDel(ctx, key, fields...)
}

// GetTranslationKey 获取翻译缓存键
func (s *CacheService) GetTranslationKey(projectID uint) string {
	return fmt.Sprintf("%s%d", domain.TranslationKeyPrefix, projectID)
}

// GetTranslationMatrixKey 获取翻译矩阵缓存键
func (s *CacheService) GetTranslationMatrixKey(projectID uint, keyword string) string {
	if keyword == "" {
		return fmt.Sprintf("%s%d", domain.TranslationMatrixPrefix, projectID)
	}
	return fmt.Sprintf("%s%d:%s", domain.TranslationMatrixPrefix, projectID, keyword)
}

// GetDashboardStatsKey 获取仪表板统计缓存键
func (s *CacheService) GetDashboardStatsKey() string {
	return domain.DashboardStatsKey
}

// GetLanguagesKey 获取语言列表缓存键
func (s *CacheService) GetLanguagesKey() string {
	return domain.LanguagesKey
}

// GetProjectKey 获取项目缓存键
func (s *CacheService) GetProjectKey(projectID uint) string {
	return fmt.Sprintf("%s%d", domain.ProjectKeyPrefix, projectID)
}

// GetProjectsKey 获取项目列表缓存键
func (s *CacheService) GetProjectsKey() string {
	return domain.ProjectsKey
}
