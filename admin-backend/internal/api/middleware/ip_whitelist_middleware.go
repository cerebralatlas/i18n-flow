package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// IPWhitelistMiddleware 创建IP白名单中间件
func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果白名单为空，允许所有访问
		if len(allowedIPs) == 0 {
			c.Next()
			return
		}

		// 获取客户端IP
		clientIP := getClientIPAddress(c)

		// 检查IP是否在白名单中
		if !isIPAllowed(clientIP, allowedIPs) {
			c.AbortWithStatusJSON(403, gin.H{
				"error":   "Forbidden",
				"message": "Your IP is not allowed to access this resource",
			})
			return
		}

		c.Next()
	}
}

// getClientIPAddress 获取客户端真实IP
func getClientIPAddress(c *gin.Context) string {
	// 尝试从X-Forwarded-For头获取
	forwardedFor := c.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从X-Real-IP头获取
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 使用RemoteAddr
	return c.ClientIP()
}

// isIPAllowed 检查IP是否在白名单中
func isIPAllowed(clientIP string, allowedIPs []string) bool {
	// 解析客户端IP
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	// 检查是否在白名单中
	for _, allowedIP := range allowedIPs {
		// 检查是否是CIDR格式
		if strings.Contains(allowedIP, "/") {
			_, ipNet, err := net.ParseCIDR(allowedIP)
			if err == nil && ipNet.Contains(ip) {
				return true
			}
		} else {
			// 直接比较IP
			if allowedIP == clientIP {
				return true
			}

			// 检查是否是IPv6的localhost
			if (allowedIP == "::1" && clientIP == "::1") ||
				(allowedIP == "127.0.0.1" && clientIP == "127.0.0.1") {
				return true
			}
		}
	}

	return false
}
