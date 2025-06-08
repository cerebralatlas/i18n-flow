package middleware

import (
	"i18n-flow/errors"
	"i18n-flow/service/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT鉴权中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.AbortWithError(c, errors.NewInvalidToken("未提供Authorization头"))
			return
		}

		// Bearer token格式检查
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			errors.AbortWithError(c, errors.NewInvalidToken("Authorization格式错误，应为'Bearer token'"))
			return
		}

		// 验证token
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				errors.AbortWithError(c, errors.NewAppError(errors.TokenExpired).WithError(err))
			} else {
				errors.AbortWithError(c, errors.NewInvalidToken(err.Error()))
			}
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
