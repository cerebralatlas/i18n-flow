package controller

import (
	"i18n-flow/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: service.UserService{},
	}
}

// Login 登录
// @Summary      用户登录
// @Description  使用用户名和密码获取访问令牌
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        credentials  body      service.LoginRequest  true  "登录凭证"
// @Success      200          {object}  service.LoginResponse
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req service.LoginRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用登录服务
	resp, err := c.userService.Login(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RefreshToken 刷新token
// @Summary      刷新访问令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        refresh_token  body      service.RefreshRequest  true  "刷新令牌"
// @Success      200            {object}  service.LoginResponse
// @Failure      400            {object}  map[string]string
// @Failure      401            {object}  map[string]string
// @Router       /refresh [post]
func (c *UserController) RefreshToken(ctx *gin.Context) {
	var req service.RefreshRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用刷新服务
	resp, err := c.userService.RefreshToken(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUserInfo 获取用户信息
// @Summary      获取当前用户信息
// @Description  获取已登录用户的详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Security     BearerAuth
// @Router       /user/info [get]
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户信息
	userID, _ := ctx.Get("userID")
	username, _ := ctx.Get("username")

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
	})
}
