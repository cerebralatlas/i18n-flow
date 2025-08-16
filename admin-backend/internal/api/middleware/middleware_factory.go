package middleware

import (
	"i18n-flow/internal/domain"

	"github.com/gin-gonic/gin"
)

// MiddlewareFactory 中间件工厂
// 负责管理需要依赖注入的中间件
type MiddlewareFactory struct {
	authService domain.AuthService
}

// NewMiddlewareFactory 创建中间件工厂
func NewMiddlewareFactory(authService domain.AuthService) *MiddlewareFactory {
	return &MiddlewareFactory{
		authService: authService,
	}
}

// JWTAuthMiddleware 返回配置好的JWT认证中间件
func (f *MiddlewareFactory) JWTAuthMiddleware() gin.HandlerFunc {
	return JWTAuthMiddleware(f.authService)
}

// 可以在这里添加其他需要依赖注入的中间件工厂方法
// 例如：
// func (f *MiddlewareFactory) RateLimitMiddleware() gin.HandlerFunc {
//     return RateLimitMiddleware(f.rateLimiter)
// }
