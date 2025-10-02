package middleware

import (
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
)

// TollboothRateLimitMiddleware 使用 tollbooth 的限流中间件
func TollboothRateLimitMiddleware(max float64, ttl time.Duration) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(max, &limiter.ExpirableOptions{
		DefaultExpirationTTL: ttl,
	})

	// 设置自定义错误消息
	lmt.SetMessage("请求过于频繁，请稍后再试")
	lmt.SetMessageContentType("application/json")

	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(429, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "请求过于频繁，请稍后再试",
					"details": "Rate limit exceeded",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 预定义的限流中间件
func TollboothGlobalRateLimit() gin.HandlerFunc {
	// 每秒100个请求，5分钟过期
	return TollboothRateLimitMiddleware(100, 5*time.Minute)
}

func TollboothLoginRateLimit() gin.HandlerFunc {
	// 每秒5个请求，10分钟过期（防暴力破解）
	return TollboothRateLimitMiddleware(5, 10*time.Minute)
}

func TollboothAPIRateLimit() gin.HandlerFunc {
	// 每秒50个请求，5分钟过期
	return TollboothRateLimitMiddleware(50, 5*time.Minute)
}

func TollboothBatchRateLimit() gin.HandlerFunc {
	// 每秒2个请求，10分钟过期
	return TollboothRateLimitMiddleware(2, 10*time.Minute)
}
