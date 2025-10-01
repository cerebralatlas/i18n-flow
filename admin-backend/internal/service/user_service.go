package service

import (
	"context"
	"i18n-flow/internal/domain"
	"strings"

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

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, domain.ErrUserExists
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
			return nil, domain.ErrEmailExists
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Status:   "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 不返回密码
	user.Password = ""
	return user, nil
}

// GetAllUsers 获取用户列表
func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int, keyword string) ([]*domain.User, int64, error) {
	users, total, err := s.userRepo.GetAll(ctx, limit, offset, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 清除密码字段
	for _, user := range users {
		user.Password = ""
	}

	return users, total, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 不返回密码
	user.Password = ""
	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, id uint, req domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Username != "" && req.Username != user.Username {
		// 检查新用户名是否已存在
		if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
			return nil, domain.ErrUserExists
		}
		user.Username = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		// 检查新邮箱是否已存在
		if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
			return nil, domain.ErrEmailExists
		}
		user.Email = req.Email
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// 不返回密码
	user.Password = ""
	return user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID uint, req domain.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return domain.ErrInvalidPassword
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// ResetPassword 重置用户密码（管理员功能）
func (s *UserService) ResetPassword(ctx context.Context, userID uint, req domain.ResetPasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 不能删除管理员用户
	if strings.ToLower(user.Role) == "admin" {
		return domain.ErrCannotDeleteAdmin
	}

	return s.userRepo.Delete(ctx, id)
}
