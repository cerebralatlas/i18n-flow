package service

import (
	"context"
	"fmt"
	"i18n-flow/internal/domain"
	"sync"
)

// CachedUserService 带缓存的用户服务实现
type CachedUserService struct {
	userService  *UserService
	cacheService domain.CacheService
	// 用于防止缓存击穿的互斥锁
	cacheMutexes map[string]*sync.Mutex
	mutexLock    sync.RWMutex
}

// NewCachedUserService 创建带缓存的用户服务实例
func NewCachedUserService(
	userService *UserService,
	cacheService domain.CacheService,
) *CachedUserService {
	return &CachedUserService{
		userService:  userService,
		cacheService: cacheService,
		cacheMutexes: make(map[string]*sync.Mutex),
	}
}

// getMutex 获取指定键的互斥锁，用于防止缓存击穿
func (s *CachedUserService) getMutex(key string) *sync.Mutex {
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
func (s *CachedUserService) removeMutex(key string) {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()

	delete(s.cacheMutexes, key)
}

// Login 用户登录
func (s *CachedUserService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// 登录操作不缓存，直接调用基础服务
	return s.userService.Login(ctx, req)
}

// RefreshToken 刷新token
func (s *CachedUserService) RefreshToken(ctx context.Context, req domain.RefreshRequest) (*domain.LoginResponse, error) {
	// 刷新token操作不缓存，直接调用基础服务
	return s.userService.RefreshToken(ctx, req)
}

// GetUserInfo 获取用户信息（使用缓存）
func (s *CachedUserService) GetUserInfo(ctx context.Context, userID uint) (*domain.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)

	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	var user *domain.User
	err := s.cacheService.GetJSONWithEmptyCheck(ctx, cacheKey, &user)
	if err == nil {
		return user, nil
	}

	// 缓存未命中，从数据库获取
	user, err = s.userService.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 更新缓存，添加随机过期时间防止雪崩
	expiration := s.cacheService.AddRandomExpiration(domain.DefaultExpiration)
	if err := s.cacheService.SetJSONWithEmptyCache(ctx, cacheKey, user, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果
	}

	return user, nil
}
