package handlers

import (
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"
	"strconv"

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
		switch err {
		case domain.ErrUserNotFound, domain.ErrInvalidPassword:
			response.Unauthorized(ctx, err.Error())
		default:
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
		switch err {
		case domain.ErrInvalidToken:
			response.InvalidToken(ctx, err.Error())
		default:
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

// CreateUser 创建用户
// @Summary      创建新用户
// @Description  管理员创建新用户账户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        user  body      domain.CreateUserRequest  true  "用户信息"
// @Success      201   {object}  domain.User
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Security     BearerAuth
// @Router       /users [post]
func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req domain.CreateUserRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 调用创建用户服务
	user, err := h.userService.CreateUser(ctx.Request.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrUserExists:
			response.Conflict(ctx, "用户名已存在")
		case domain.ErrEmailExists:
			response.Conflict(ctx, "邮箱已存在")
		default:
			response.InternalServerError(ctx, "创建用户失败")
		}
		return
	}

	response.Created(ctx, user)
}

// GetUsers 获取用户列表
// @Summary      获取用户列表
// @Description  分页获取用户列表，支持关键词搜索
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "页码"        default(1)
// @Param        page_size query     int     false  "每页数量"     default(10)
// @Param        keyword   query     string  false  "搜索关键词"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]string
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) GetUsers(ctx *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	keyword := ctx.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取用户列表
	users, total, err := h.userService.GetAllUsers(ctx.Request.Context(), pageSize, offset, keyword)
	if err != nil {
		response.InternalServerError(ctx, "获取用户列表失败")
		return
	}

	meta := &response.Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}

	response.SuccessWithMeta(ctx, users, meta)
}

// GetUser 获取用户详情
// @Summary      获取用户详情
// @Description  根据用户ID获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  domain.User
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(ctx *gin.Context) {
	// 解析用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ValidationError(ctx, "无效的用户ID")
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(ctx.Request.Context(), uint(id))
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			response.NotFound(ctx, "用户不存在")
		default:
			response.InternalServerError(ctx, "获取用户信息失败")
		}
		return
	}

	response.Success(ctx, user)
}

// UpdateUser 更新用户
// @Summary      更新用户信息
// @Description  更新用户的基本信息和角色状态
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "用户ID"
// @Param        user  body      domain.UpdateUserRequest  true  "用户信息"
// @Success      200   {object}  domain.User
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	// 解析用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ValidationError(ctx, "无效的用户ID")
		return
	}

	var req domain.UpdateUserRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 调用更新用户服务
	user, err := h.userService.UpdateUser(ctx.Request.Context(), uint(id), req)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			response.NotFound(ctx, "用户不存在")
		case domain.ErrUserExists:
			response.Conflict(ctx, "用户名已存在")
		case domain.ErrEmailExists:
			response.Conflict(ctx, "邮箱已存在")
		default:
			response.InternalServerError(ctx, "更新用户失败")
		}
		return
	}

	response.Success(ctx, user)
}

// ChangePassword 修改密码
// @Summary      修改用户密码
// @Description  用户修改自己的密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        password  body      domain.ChangePasswordRequest  true  "密码信息"
// @Success      200       {object}  map[string]string
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Security     BearerAuth
// @Router       /user/change-password [post]
func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	var req domain.ChangePasswordRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 调用修改密码服务
	if err := h.userService.ChangePassword(ctx.Request.Context(), userID.(uint), req); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			response.NotFound(ctx, "用户不存在")
		case domain.ErrInvalidPassword:
			response.Unauthorized(ctx, "原密码错误")
		default:
			response.InternalServerError(ctx, "修改密码失败")
		}
		return
	}

	response.Success(ctx, map[string]string{"message": "密码修改成功"})
}

// ResetPassword 重置用户密码
// @Summary      重置用户密码
// @Description  管理员重置指定用户的密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id        path      int                       true  "用户ID"
// @Param        password  body      domain.ResetPasswordRequest  true  "新密码信息"
// @Success      200       {object}  map[string]string
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Security     BearerAuth
// @Router       /users/{id}/reset-password [post]
func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	// 解析用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ValidationError(ctx, "无效的用户ID")
		return
	}

	var req domain.ResetPasswordRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 调用重置密码服务
	if err := h.userService.ResetPassword(ctx.Request.Context(), uint(id), req); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			response.NotFound(ctx, "用户不存在")
		default:
			response.InternalServerError(ctx, "重置密码失败")
		}
		return
	}

	response.Success(ctx, map[string]string{"message": "密码重置成功"})
}

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  删除指定的用户账户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	// 解析用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ValidationError(ctx, "无效的用户ID")
		return
	}

	// 调用删除用户服务
	if err := h.userService.DeleteUser(ctx.Request.Context(), uint(id)); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			response.NotFound(ctx, "用户不存在")
		case domain.ErrCannotDeleteAdmin:
			response.Forbidden(ctx, "不能删除管理员用户")
		default:
			response.InternalServerError(ctx, "删除用户失败")
		}
		return
	}

	response.Success(ctx, map[string]string{"message": "用户删除成功"})
}
