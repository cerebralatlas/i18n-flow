package handlers

import (
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService domain.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login 登录
// @Summary      用户登录
// @Description  使用用户名和密码获取访问令牌
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        credentials  body      domain.LoginRequest  true  "登录凭证"
// @Success      200          {object}  domain.LoginResponse
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req domain.LoginRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 调用登录服务
	resp, err := h.userService.Login(ctx.Request.Context(), req)
	if err != nil {
		// 根据错误类型返回不同状态码
		if err == domain.ErrUserNotFound || err == domain.ErrInvalidPassword {
			response.Unauthorized(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "登录失败")
		}
		return
	}

	response.Success(ctx, resp)
}

// RefreshToken 刷新token
// @Summary      刷新访问令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        refresh_token  body      domain.RefreshRequest  true  "刷新令牌"
// @Success      200            {object}  domain.LoginResponse
// @Failure      400            {object}  map[string]string
// @Failure      401            {object}  map[string]string
// @Router       /refresh [post]
func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	var req domain.RefreshRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 调用刷新服务
	resp, err := h.userService.RefreshToken(ctx.Request.Context(), req)
	if err != nil {
		if err == domain.ErrInvalidToken {
			response.InvalidToken(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "刷新token失败")
		}
		return
	}

	response.Success(ctx, resp)
}

// GetUserInfo 获取用户信息
// @Summary      获取当前用户信息
// @Description  获取已登录用户的详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.User
// @Failure      401  {object}  map[string]string
// @Security     BearerAuth
// @Router       /user/info [get]
func (h *UserHandler) GetUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserInfo(ctx.Request.Context(), userID.(uint))
	if err != nil {
		response.InternalServerError(ctx, "获取用户信息失败")
		return
	}

	response.Success(ctx, user)
}
