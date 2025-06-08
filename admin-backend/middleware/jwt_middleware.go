package middleware

import (
	"i18n-flow/errors"
	"i18n-flow/service/auth"
	"i18n-flow/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
			// 记录认证失败日志
			utils.AuthLog("JWT authentication failed",
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
			)

			if strings.Contains(err.Error(), "expired") {
				errors.AbortWithError(c, errors.NewAppError(errors.TokenExpired).WithError(err))
			} else {
				errors.AbortWithError(c, errors.NewInvalidToken(err.Error()))
			}
			return
		}

		// 记录认证成功日志
		utils.AuthLog("JWT authentication successful",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
