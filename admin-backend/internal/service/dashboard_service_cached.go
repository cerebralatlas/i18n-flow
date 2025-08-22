package service

import (
	"context"
	"fmt"
	"i18n-flow/internal/domain"
)

// CachedDashboardService 带缓存的仪表板服务实现
type CachedDashboardService struct {
	dashboardService *DashboardService
	cacheService     domain.CacheService
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

// GetStats 获取仪表板统计信息（使用缓存）
func (s *CachedDashboardService) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	// 尝试从缓存获取
	var stats domain.DashboardStats
	err := s.cacheService.GetJSON(ctx, s.cacheService.GetDashboardStatsKey(), &stats)
	if err == nil {
		return &stats, nil
	}

	// 缓存未命中，从数据库获取
	dbStats, err := s.dashboardService.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	stats = *dbStats

	// 更新缓存
	if err := s.cacheService.SetJSON(ctx, s.cacheService.GetDashboardStatsKey(), stats, domain.LongExpiration); err != nil {
		// 缓存更新失败，但不影响返回结果
		fmt.Printf("缓存更新失败: %v\n", err)
	}

	return &stats, nil
}
