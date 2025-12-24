package domain

import (
	"context"
	"time"
)

// CacheService 缓存服务接口
type CacheService interface {
	// 基础操作
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)

	// JSON操作
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error

	// 带空值缓存的操作
	SetWithEmptyCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetWithEmptyCheck(ctx context.Context, key string) (string, error)
	SetJSONWithEmptyCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSONWithEmptyCheck(ctx context.Context, key string, dest interface{}) error

	// 哈希表操作
	HSet(ctx context.Context, key, field string, value interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// 缓存键生成
	GetTranslationKey(projectID uint64) string
	GetTranslationMatrixKey(projectID uint64, keyword string) string
	GetDashboardStatsKey() string
	GetLanguagesKey() string
	GetProjectKey(projectID uint64) string
	GetProjectsKey() string
	
	// 添加随机过期时间防止雪崩
	AddRandomExpiration(baseExpiration time.Duration) time.Duration
}

// 缓存键常量
const (
	// 过期时间
	DefaultExpiration = 30 * time.Minute
	LongExpiration    = 12 * time.Hour
	ShortExpiration   = 5 * time.Minute // 用于空值缓存

	// 缓存键前缀
	TranslationKeyPrefix    = "translation:"
	TranslationMatrixPrefix = "translation_matrix:"
	DashboardStatsKey       = "dashboard:stats"
	LanguagesKey            = "languages"
	ProjectKeyPrefix        = "project:"
	ProjectsKey             = "projects"
)

// ErrCacheMiss 缓存未命中错误
var ErrCacheMiss = CacheError("cache miss")

// CacheError 缓存错误类型
type CacheError string

func (e CacheError) Error() string {
	return string(e)
}
