package service

import (
	"context"
	"i18n-flow/internal/domain"
	"i18n-flow/internal/dto"
	"sync"
)

// CachedDashboardService 带缓存的仪表板服务实现
type CachedDashboardService struct {
	dashboardService *DashboardService
	cacheService     domain.CacheService
	// 用于防止缓存击穿的互斥锁，使用 sync.Map 线程安全
	cacheMutexes sync.Map
}

// NewCachedDashboardService 创建带缓存的仪表板服务实例
func NewCachedDashboardService(
	dashboardService *DashboardService,
	cacheService domain.CacheService,
) *CachedDashboardService {
	return &CachedDashboardService{
		dashboardService: dashboardService,
		cacheService:     cacheService,
	}
}

// getMutex 获取指定键的互斥锁，用于防止缓存击穿
func (s *CachedDashboardService) getMutex(key string) *sync.Mutex {
	if mutex, exists := s.cacheMutexes.Load(key); exists {
		return mutex.(*sync.Mutex)
	}

	mutex := &sync.Mutex{}
	actual, loaded := s.cacheMutexes.LoadOrStore(key, mutex)
	if loaded {
		return actual.(*sync.Mutex)
	}
	return mutex
}

// removeMutex 移除指定键的互斥锁
func (s *CachedDashboardService) removeMutex(key string) {
	s.cacheMutexes.Delete(key)
}

// GetStats 获取仪表板统计信息（使用缓存）
func (s *CachedDashboardService) GetStats(ctx context.Context) (*dto.DashboardStats, error) {
	cacheKey := s.cacheService.GetDashboardStatsKey()

	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	var stats *dto.DashboardStats
	err := s.cacheService.GetJSONWithEmptyCheck(ctx, cacheKey, &stats)
	if err == nil {
		return stats, nil
	}

	// 缓存未命中，从数据库获取
	stats, err = s.dashboardService.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	// 更新缓存，添加随机过期时间防止雪崩
	expiration := s.cacheService.AddRandomExpiration(domain.LongExpiration)
	if err := s.cacheService.SetJSONWithEmptyCache(ctx, cacheKey, stats, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果
		//fmt.Printf("缓存更新失败: %v\n", err)
	}

	return stats, nil
}
