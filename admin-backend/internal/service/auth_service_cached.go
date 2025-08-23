package service

import (
	"context"
	"i18n-flow/internal/domain"
	"sync"
	"time"
)

// CachedAuthService 带缓存的认证服务实现
type CachedAuthService struct {
	authService *AuthService
	cacheService domain.CacheService
	// 用于防止缓存击穿的互斥锁
	cacheMutexes map[string]*sync.Mutex
	mutexLock    sync.RWMutex
}

// NewCachedAuthService 创建带缓存的认证服务实例
func NewCachedAuthService(
	authService *AuthService,
	cacheService domain.CacheService,
) *CachedAuthService {
	return &CachedAuthService{
		authService:  authService,
		cacheService: cacheService,
		cacheMutexes: make(map[string]*sync.Mutex),
	}
}

// getMutex 获取指定键的互斥锁，用于防止缓存击穿
func (s *CachedAuthService) getMutex(key string) *sync.Mutex {
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
func (s *CachedAuthService) removeMutex(key string) {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()
	
	delete(s.cacheMutexes, key)
}

// GenerateToken 生成JWT token
func (s *CachedAuthService) GenerateToken(user *domain.User) (string, error) {
	// 生成token操作不缓存，直接调用基础服务
	return s.authService.GenerateToken(user)
}

// GenerateRefreshToken 生成刷新token
func (s *CachedAuthService) GenerateRefreshToken(user *domain.User) (string, error) {
	// 生成刷新token操作不缓存，直接调用基础服务
	return s.authService.GenerateRefreshToken(user)
}

// ValidateToken 验证JWT token（使用缓存）
func (s *CachedAuthService) ValidateToken(token string) (*domain.User, error) {
	cacheKey := "token:" + token
	
	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	var user *domain.User
	err := s.cacheService.GetJSONWithEmptyCheck(context.Background(), cacheKey, &user)
	if err == nil {
		return user, nil
	}

	// 缓存未命中，从数据库获取
	user, err = s.authService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// 更新缓存，设置较短的过期时间（token有效期通常较短）
	expiration := s.cacheService.AddRandomExpiration(5 * time.Minute)
	if err := s.cacheService.SetJSONWithEmptyCache(context.Background(), cacheKey, user, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果
	}

	return user, nil
}

// ValidateRefreshToken 验证刷新token（使用缓存）
func (s *CachedAuthService) ValidateRefreshToken(token string) (*domain.User, error) {
	cacheKey := "refresh_token:" + token
	
	// 使用互斥锁防止缓存击穿
	mutex := s.getMutex(cacheKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.removeMutex(cacheKey) // 请求完成后移除锁
	}()

	// 尝试从缓存获取
	var user *domain.User
	err := s.cacheService.GetJSONWithEmptyCheck(context.Background(), cacheKey, &user)
	if err == nil {
		return user, nil
	}

	// 缓存未命中，从数据库获取
	user, err = s.authService.ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	// 更新缓存，设置较短的过期时间（refresh token有效期通常较短）
	expiration := s.cacheService.AddRandomExpiration(30 * time.Minute)
	if err := s.cacheService.SetJSONWithEmptyCache(context.Background(), cacheKey, user, expiration); err != nil {
		// 缓存更新失败，但不影响返回结果
	}

	return user, nil
}