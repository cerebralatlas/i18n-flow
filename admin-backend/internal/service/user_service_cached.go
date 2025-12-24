package service

import (
	"context"
	"fmt"
	"i18n-flow/internal/domain"
	"sync"
	"time"
)

// CachedUserService 带缓存的用户服务实现
type CachedUserService struct {
	userService  *UserService
	cacheService domain.CacheService
	// 用于防止缓存击穿的互斥锁，使用 sync.Map 线程安全
	cacheMutexes sync.Map
}

// NewCachedUserService 创建带缓存的用户服务实例
func NewCachedUserService(
	userService *UserService,
	cacheService domain.CacheService,
) *CachedUserService {
	svc := &CachedUserService{
		userService:  userService,
		cacheService: cacheService,
	}
	// 启动清理协程
	go svc.cleanupMutexes()
	return svc
}

// getMutex 获取指定键的互斥锁，用于防止缓存击穿
func (s *CachedUserService) getMutex(key string) *sync.Mutex {
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
func (s *CachedUserService) removeMutex(key string) {
	s.cacheMutexes.Delete(key)
}

// cleanupMutexes 定期清理无效的 mutex 锁
func (s *CachedUserService) cleanupMutexes() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// 由于每次请求后都会调用 removeMutex，map 不会无限增长
	}
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
func (s *CachedUserService) GetUserInfo(ctx context.Context, userID uint64) (*domain.User, error) {
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

// CreateUser 创建用户（不缓存）
func (s *CachedUserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	return s.userService.CreateUser(ctx, req)
}

// GetAllUsers 获取用户列表（不缓存）
func (s *CachedUserService) GetAllUsers(ctx context.Context, limit, offset int, keyword string) ([]*domain.User, int64, error) {
	return s.userService.GetAllUsers(ctx, limit, offset, keyword)
}

// GetUserByID 根据ID获取用户（使用缓存）
func (s *CachedUserService) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	// 复用GetUserInfo的缓存逻辑
	return s.GetUserInfo(ctx, id)
}

// UpdateUser 更新用户（清除缓存）
func (s *CachedUserService) UpdateUser(ctx context.Context, id uint64, req domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userService.UpdateUser(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// 清除用户缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	s.cacheService.Delete(ctx, cacheKey)

	return user, nil
}

// ChangePassword 修改密码（不缓存）
func (s *CachedUserService) ChangePassword(ctx context.Context, userID uint64, req domain.ChangePasswordRequest) error {
	return s.userService.ChangePassword(ctx, userID, req)
}

// ResetPassword 重置密码（不缓存）
func (s *CachedUserService) ResetPassword(ctx context.Context, userID uint64, req domain.ResetPasswordRequest) error {
	return s.userService.ResetPassword(ctx, userID, req)
}

// DeleteUser 删除用户（清除缓存）
func (s *CachedUserService) DeleteUser(ctx context.Context, id uint64) error {
	err := s.userService.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// 清除用户缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	s.cacheService.Delete(ctx, cacheKey)

	return nil
}
