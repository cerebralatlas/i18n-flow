package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware 安全HTTP头中间件
// 设置各种安全相关的HTTP响应头，防护常见的Web安全攻击
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止内容类型嗅探攻击
		c.Header("X-Content-Type-Options", "nosniff")

		// 防止点击劫持攻击
		c.Header("X-Frame-Options", "DENY")

		// XSS保护（虽然现代浏览器默认启用，但为了兼容性还是设置）
		c.Header("X-XSS-Protection", "1; mode=block")

		// 强制HTTPS传输安全（生产环境建议启用）
		// 注意：仅在HTTPS环境下设置此头
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// 内容安全策略 - 防止XSS攻击
		// 更严格的CSP策略，移除unsafe-inline和unsafe-eval，添加违规报告
		csp := "default-src 'self'; " +
			"script-src 'self'; " +
			"style-src 'self'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'; " +
			"object-src 'none'; " +
			"media-src 'none'; " +
			"worker-src 'none'; " +
			"child-src 'none'; " +
			"frame-src 'none'; " +
			"report-uri /csp-report"
		c.Header("Content-Security-Policy", csp)

		// 引用者策略 - 控制Referer头的发送
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 权限策略 - 禁用不需要的浏览器功能
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), interest-cohort=()")

		// 移除可能泄露服务器信息的头
		c.Header("Server", "")
		c.Header("X-Powered-By", "")

		// 防止缓存敏感信息
		if c.Request.URL.Path == "/api/login" ||
			c.Request.URL.Path == "/api/refresh" ||
			c.Request.URL.Path == "/api/user/info" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableHSTS     bool     // 是否启用HSTS
	HSTSMaxAge     int      // HSTS最大年龄（秒）
	EnableCSP      bool     // 是否启用CSP
	CustomCSP      string   // 自定义CSP策略
	AllowedOrigins []string // 允许的源
	TrustedProxies []string // 受信任的代理
}

// SecurityMiddlewareWithConfig 带配置的安全中间件
func SecurityMiddlewareWithConfig(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 基础安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Server", "")
		c.Header("X-Powered-By", "")

		// 条件性HSTS
		if config.EnableHSTS && (c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https") {
			maxAge := config.HSTSMaxAge
			if maxAge == 0 {
				maxAge = 31536000 // 默认1年
			}
			c.Header("Strict-Transport-Security", fmt.Sprintf("max-age=%d; includeSubDomains; preload", maxAge))
		}

		// 条件性CSP
		if config.EnableCSP {
			csp := config.CustomCSP
			if csp == "" {
				csp = "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"
			}
			c.Header("Content-Security-Policy", csp)
		}

		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
}
