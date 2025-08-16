package middleware

import (
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT鉴权中间件
// 接受authService作为参数，支持依赖注入
func JWTAuthMiddleware(authService domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供Authorization头")
			return
		}

		// Bearer token格式检查
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.BadRequest(c, "Authorization格式错误，应为'Bearer token'")
			return
		}

		// 验证token
		tokenString := parts[1]
		user, err := authService.ValidateToken(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				response.TokenExpired(c, "token已过期")
			} else {
				response.InvalidToken(c, "无效的token")
			}
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", user.ID)
		c.Set("username", user.Username)

		c.Next()
	}
}
