package middleware

import (
	"fmt"
	"i18n-flow/internal/api/response"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitManager 限流管理器
type RateLimitManager struct {
	visitors map[string]*RateLimiter
	mutex    sync.RWMutex
	rate     rate.Limit // 请求速率（每秒请求数）
	burst    int        // 突发请求数
}

// NewRateLimitManager 创建限流管理器
func NewRateLimitManager(requestsPerSecond float64, burst int) *RateLimitManager {
	manager := &RateLimitManager{
		visitors: make(map[string]*RateLimiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}

	// 启动清理协程，定期清理过期的限流器
	go manager.cleanupVisitors()
	return manager
}

// getRateLimiter 获取指定IP的限流器
func (rlm *RateLimitManager) getRateLimiter(ip string) *rate.Limiter {
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()

	limiter, exists := rlm.visitors[ip]
	if !exists {
		limiter = &RateLimiter{
			limiter:  rate.NewLimiter(rlm.rate, rlm.burst),
			lastSeen: time.Now(),
		}
		rlm.visitors[ip] = limiter
	} else {
		limiter.lastSeen = time.Now()
	}

	return limiter.limiter
}

// cleanupVisitors 定期清理过期的访问者
func (rlm *RateLimitManager) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rlm.mutex.Lock()
		for ip, limiter := range rlm.visitors {
			// 清理5分钟内未访问的限流器
			if time.Since(limiter.lastSeen) > 5*time.Minute {
				delete(rlm.visitors, ip)
			}
		}
		rlm.mutex.Unlock()
	}
}

// 全局限流管理器实例
var (
	globalRateLimiter   *RateLimitManager
	loginRateLimiter    *RateLimitManager
	apiRateLimiter      *RateLimitManager
	batchRateLimiter    *RateLimitManager
	rateLimiterInitOnce sync.Once
)

// initRateLimiters 初始化限流器
func initRateLimiters() {
	rateLimiterInitOnce.Do(func() {
		// 全局限流：每秒100个请求，突发200个
		globalRateLimiter = NewRateLimitManager(100, 200)

		// 登录限流：每秒5个请求，突发10个（防止暴力破解）
		loginRateLimiter = NewRateLimitManager(5, 10)

		// API限流：每秒50个请求，突发100个
		apiRateLimiter = NewRateLimitManager(50, 100)

		// 批量操作限流：每秒2个请求，突发5个
		batchRateLimiter = NewRateLimitManager(2, 5)
	})
}

// RateLimitMiddleware 通用限流中间件
func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	manager := NewRateLimitManager(requestsPerSecond, burst)

	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := manager.getRateLimiter(ip)

		if !limiter.Allow() {
			response.ErrorWithDetails(c, 429, "RATE_LIMIT_EXCEEDED",
				"请求过于频繁，请稍后再试",
				fmt.Sprintf("IP: %s 超出限流限制", ip))
			return
		}

		c.Next()
	}
}

// GlobalRateLimitMiddleware 全局限流中间件
func GlobalRateLimitMiddleware() gin.HandlerFunc {
	initRateLimiters()

	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := globalRateLimiter.getRateLimiter(ip)

		if !limiter.Allow() {
			response.ErrorWithDetails(c, 429, "RATE_LIMIT_EXCEEDED",
				"请求过于频繁，请稍后再试",
				fmt.Sprintf("全局限流: IP %s", ip))
			return
		}

		c.Next()
	}
}

// LoginRateLimitMiddleware 登录限流中间件
func LoginRateLimitMiddleware() gin.HandlerFunc {
	initRateLimiters()

	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := loginRateLimiter.getRateLimiter(ip)

		if !limiter.Allow() {
			response.ErrorWithDetails(c, 429, "LOGIN_RATE_LIMIT_EXCEEDED",
				"登录尝试过于频繁，请稍后再试",
				fmt.Sprintf("登录限流: IP %s", ip))
			return
		}

		c.Next()
	}
}

// APIRateLimitMiddleware API限流中间件
func APIRateLimitMiddleware() gin.HandlerFunc {
	initRateLimiters()

	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := apiRateLimiter.getRateLimiter(ip)

		if !limiter.Allow() {
			response.ErrorWithDetails(c, 429, "API_RATE_LIMIT_EXCEEDED",
				"API请求过于频繁，请稍后再试",
				fmt.Sprintf("API限流: IP %s", ip))
			return
		}

		c.Next()
	}
}

// BatchOperationRateLimitMiddleware 批量操作限流中间件
func BatchOperationRateLimitMiddleware() gin.HandlerFunc {
	initRateLimiters()

	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := batchRateLimiter.getRateLimiter(ip)

		if !limiter.Allow() {
			response.ErrorWithDetails(c, 429, "BATCH_RATE_LIMIT_EXCEEDED",
				"批量操作过于频繁，请稍后再试",
				fmt.Sprintf("批量操作限流: IP %s", ip))
			return
		}

		c.Next()
	}
}

// getClientIP 获取客户端真实IP
func getClientIP(c *gin.Context) string {
	// 优先检查X-Real-IP头
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// 检查X-Forwarded-For头（可能包含多个IP）
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		// 取第一个IP
		if firstIP := net.ParseIP(ip); firstIP != nil {
			return ip
		}
	}

	// 使用Gin的ClientIP方法作为后备
	return c.ClientIP()
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	GlobalRate  float64  // 全局每秒请求数
	GlobalBurst int      // 全局突发请求数
	LoginRate   float64  // 登录每秒请求数
	LoginBurst  int      // 登录突发请求数
	APIRate     float64  // API每秒请求数
	APIBurst    int      // API突发请求数
	BatchRate   float64  // 批量操作每秒请求数
	BatchBurst  int      // 批量操作突发请求数
	TrustedIPs  []string // 受信任的IP列表（不限流）
}

// RateLimitMiddlewareWithConfig 带配置的限流中间件
func RateLimitMiddlewareWithConfig(config RateLimitConfig) gin.HandlerFunc {
	manager := NewRateLimitManager(config.GlobalRate, config.GlobalBurst)

	// 预编译受信任的IP列表
	trustedIPNets := make([]*net.IPNet, 0, len(config.TrustedIPs))
	for _, ipStr := range config.TrustedIPs {
		if ip := net.ParseIP(ipStr); ip != nil {
			ipNet := &net.IPNet{
				IP:   ip,
				Mask: net.CIDRMask(32, 32), // 默认/32掩码
			}
			trustedIPNets = append(trustedIPNets, ipNet)
		} else if _, ipNet, err := net.ParseCIDR(ipStr); err == nil {
			trustedIPNets = append(trustedIPNets, ipNet)
		}
	}

	return func(c *gin.Context) {
		ip := getClientIP(c)

		// 检查是否为受信任IP
		clientIP := net.ParseIP(ip)
		if clientIP != nil {
			for _, trustedNet := range trustedIPNets {
				if trustedNet.Contains(clientIP) {
					c.Next()
					return
				}
			}
		}

		limiter := manager.getRateLimiter(ip)
		if !limiter.Allow() {
			response.ErrorWithDetails(c, 429, "RATE_LIMIT_EXCEEDED",
				"请求过于频繁，请稍后再试",
				fmt.Sprintf("IP: %s", ip))
			return
		}

		c.Next()
	}
}
