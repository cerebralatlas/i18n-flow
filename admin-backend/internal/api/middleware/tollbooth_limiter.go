package middleware

import (
	"fmt"
	"i18n-flow/internal/api/response"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
)

// TollboothLimitMiddleware 使用 tollbooth 的通用限流中间件
func TollboothLimitMiddleware(max float64, ttl time.Duration, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(max, &limiter.ExpirableOptions{
		DefaultExpirationTTL: ttl,
	})

	// 设置 IP 提取函数
	if keyFunc != nil {
		lmt.SetIPLookups([]string{"X-Real-IP", "X-Forwarded-For", "RemoteAddr"})
	}

	return func(c *gin.Context) {
		// 获取限流键（通常是 IP）
		var key string
		if keyFunc != nil {
			key = keyFunc(c)
		} else {
			key = getClientIP(c)
		}

		// 检查限流
		err := tollbooth.LimitByKeys(lmt, []string{key})
		if err != nil {
			response.ErrorWithDetails(c, 429, "RATE_LIMIT_EXCEEDED",
				"请求过于频繁，请稍后再试",
				fmt.Sprintf("Rate limit exceeded for: %s", key))
			c.Abort()
			return
		}

		c.Next()
	}
}

// TollboothGlobalRateLimitMiddleware 全局限流中间件
func TollboothGlobalRateLimitMiddleware() gin.HandlerFunc {
	// 每秒100个请求，5分钟过期
	return TollboothLimitMiddleware(100, 5*time.Minute, nil)
}

// TollboothLoginRateLimitMiddleware 登录限流中间件
func TollboothLoginRateLimitMiddleware() gin.HandlerFunc {
	// 每秒5个请求，10分钟过期（防止暴力破解）
	return TollboothLimitMiddleware(5, 10*time.Minute, nil)
}

// TollboothAPIRateLimitMiddleware API限流中间件
func TollboothAPIRateLimitMiddleware() gin.HandlerFunc {
	// 每秒50个请求，5分钟过期
	return TollboothLimitMiddleware(50, 5*time.Minute, nil)
}

// TollboothBatchOperationRateLimitMiddleware 批量操作限流中间件
func TollboothBatchOperationRateLimitMiddleware() gin.HandlerFunc {
	// 每秒2个请求，10分钟过期
	return TollboothLimitMiddleware(2, 10*time.Minute, nil)
}

// TollboothCustomRateLimitMiddleware 自定义限流中间件
func TollboothCustomRateLimitMiddleware(max float64, ttl time.Duration) gin.HandlerFunc {
	return TollboothLimitMiddleware(max, ttl, nil)
}

// TollboothUserBasedRateLimitMiddleware 基于用户的限流中间件
func TollboothUserBasedRateLimitMiddleware(max float64, ttl time.Duration) gin.HandlerFunc {
	return TollboothLimitMiddleware(max, ttl, func(c *gin.Context) string {
		// 优先使用用户ID，如果没有则使用IP
		if userID, exists := c.Get("userID"); exists {
			return fmt.Sprintf("user:%v", userID)
		}
		return fmt.Sprintf("ip:%s", getClientIP(c))
	})
}
