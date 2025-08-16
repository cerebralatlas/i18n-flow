package service

import (
	"context"
	"i18n-flow/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务实现
type UserService struct {
	userRepo    domain.UserRepository
	authService domain.AuthService
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo domain.UserRepository, authService domain.AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
	}
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// 查询用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, domain.ErrInvalidPassword
	}

	// 生成JWT token
	token, err := s.authService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// 生成刷新token
	refreshToken, err := s.authService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// 不返回密码
	userResponse := *user
	userResponse.Password = ""

	return &domain.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         userResponse,
	}, nil
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(ctx context.Context, req domain.RefreshRequest) (*domain.LoginResponse, error) {
	// 验证刷新token
	userFromToken, err := s.authService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// 查询用户确保用户仍然存在
	user, err := s.userRepo.GetByID(ctx, userFromToken.ID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// 生成新token
	token, err := s.authService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// 生成新刷新token
	refreshToken, err := s.authService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// 不返回密码
	userResponse := *user
	userResponse.Password = ""

	return &domain.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         userResponse,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID uint) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 不返回密码
	user.Password = ""
	return user, nil
}
