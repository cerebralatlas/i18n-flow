package service

import (
	"context"
	"fmt"
	"i18n-flow/internal/domain"
	"sync"
)

// CachedProjectService 带缓存的项目服务实现
type CachedProjectService struct {
	projectService *ProjectService
	cacheService   domain.CacheService
	// 用于防止缓存击穿的互斥锁
	cacheMutexes map[string]*sync.Mutex
	mutexLock    sync.RWMutex
}

// NewCachedProjectService 创建带缓存的项目服务实例
func NewCachedProjectService(
	projectService *ProjectService,
	cacheService domain.CacheService,
) *CachedProjectService {
	return &CachedProjectService{
		projectService: projectService,
		cacheService:   cacheService,
		cacheMutexes:   make(map[string]*sync.Mutex),
	}
}

// getMutex 获取指定键的互斥锁，用于防止缓存击穿
func (s *CachedProjectService) getMutex(key string) *sync.Mutex {
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
func (s *CachedProjectService) removeMutex(key string) {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()

	delete(s.cacheMutexes, key)
}

// Create 创建项目（更新缓存）
func (s *CachedProjectService) Create(ctx context.Context, req domain.CreateProjectRequest) (*domain.Project, error) {
	project, err := s.projectService.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// 清除项目列表缓存
	s.cacheService.Delete(ctx, s.cacheService.GetProjectsKey())

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())

	return project, nil
}

// GetByID 根据ID获取项目（使用缓存）
func (s *CachedProjectService) GetByID(ctx context.Context, id uint) (*domain.Project, error) {
	cacheKey := s.cacheService.GetProjectKey(id)

	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	var project *domain.Project
	err := s.cacheService.GetJSONWithEmptyCheck(ctx, cacheKey, &project)
	if err == nil {
		return project, nil
	}

	// 缓存未命中，从数据库获取
	project, err = s.projectService.GetByID(ctx, id)
	if err != nil {
		// 对于不存在的项目，也缓存一小段时间防止缓存穿透
		if err == domain.ErrProjectNotFound {
			expiration := s.cacheService.AddRandomExpiration(domain.ShortExpiration)
			s.cacheService.SetJSONWithEmptyCache(ctx, cacheKey, nil, expiration)
		}
		return nil, err
	}

	// 更新缓存，添加随机过期时间防止雪崩
	expiration := s.cacheService.AddRandomExpiration(domain.DefaultExpiration)
	if err := s.cacheService.SetJSONWithEmptyCache(ctx, cacheKey, project, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果
	}

	return project, nil
}

// GetAll 获取所有项目（使用缓存）
func (s *CachedProjectService) GetAll(ctx context.Context, limit, offset int, keyword string) ([]*domain.Project, int64, error) {
	// 生成缓存键
	cacheKey := s.cacheService.GetProjectsKey()
	if keyword != "" {
		// 如果有搜索关键词，添加到缓存键中
		cacheKey += ":search:" + keyword
	}
	cacheKey += fmt.Sprintf(":%d:%d", limit, offset)

	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	type projectsCacheResult struct {
		Projects []*domain.Project `json:"projects"`
		Total    int64             `json:"total"`
	}

	var cachedResult projectsCacheResult
	err := s.cacheService.GetJSONWithEmptyCheck(ctx, cacheKey, &cachedResult)
	if err == nil {
		return cachedResult.Projects, cachedResult.Total, nil
	}

	// 缓存未命中，从数据库获取
	projects, total, err := s.projectService.GetAll(ctx, limit, offset, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 更新缓存，添加随机过期时间防止雪崩
	cachedResult = projectsCacheResult{
		Projects: projects,
		Total:    total,
	}

	expiration := s.cacheService.AddRandomExpiration(domain.DefaultExpiration)
	if err := s.cacheService.SetJSONWithEmptyCache(ctx, cacheKey, cachedResult, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果
	}

	return projects, total, nil
}

// Update 更新项目（更新缓存）
func (s *CachedProjectService) Update(ctx context.Context, id uint, req domain.UpdateProjectRequest) (*domain.Project, error) {
	project, err := s.projectService.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// 清除该项目的缓存
	s.cacheService.Delete(ctx, s.cacheService.GetProjectKey(id))

	// 清除项目列表缓存
	s.cacheService.Delete(ctx, s.cacheService.GetProjectsKey())

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())

	return project, nil
}

// Delete 删除项目（更新缓存）
func (s *CachedProjectService) Delete(ctx context.Context, id uint) error {
	err := s.projectService.Delete(ctx, id)
	if err != nil {
		return err
	}

	// 清除该项目的缓存
	s.cacheService.Delete(ctx, s.cacheService.GetProjectKey(id))

	// 清除项目列表缓存
	s.cacheService.Delete(ctx, s.cacheService.GetProjectsKey())

	// 清除仪表板缓存
	s.cacheService.Delete(ctx, s.cacheService.GetDashboardStatsKey())

	return nil
}
