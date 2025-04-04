package service

import (
	"errors"
	"i18n-flow/model"
	"i18n-flow/model/db"
	"i18n-flow/service/auth"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`    // 用户名
	Password string `json:"password" binding:"required" example:"password"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string     `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`         // JWT访问令牌
	RefreshToken string     `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 刷新令牌
	User         model.User `json:"user"`
}

// RefreshRequest 刷新token请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 刷新令牌
}

// UserService 用户服务
type UserService struct{}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	var user model.User

	// 查询用户
	result := db.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	// 生成刷新token
	refreshToken, err := auth.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	// 不返回密码
	user.Password = ""

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(req RefreshRequest) (*LoginResponse, error) {
	// 验证刷新token
	claims, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 查询用户
	var user model.User
	result := db.DB.First(&user, claims.UserID)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 生成新token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	// 生成新刷新token
	refreshToken, err := auth.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	// 不返回密码
	user.Password = ""

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
