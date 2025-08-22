package routes

import (
	"i18n-flow/internal/api/middleware"
	"i18n-flow/internal/container"
	"i18n-flow/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// SetupMetricsRoutes 设置监控相关路由
func SetupMetricsRoutes(router *gin.Engine, c *container.Container) {
	// 将容器注入到上下文中，以便在路由处理函数中使用
	router.Use(func(ctx *gin.Context) {
		ctx.Set("container", c)
		ctx.Next()
	})

	// 监控指标路由组
	metricsGroup := router.Group("/metrics")
	{
		// 添加IP白名单中间件，限制只有特定IP可以访问监控指标
		// 默认允许所有IP访问
		var allowedIPs []string

		// 尝试从配置中获取IP白名单
		if c.Config().Metrics.Enabled {
			allowedIPs = c.Config().Metrics.AllowedIPs
		}

		metricsGroup.Use(middleware.IPWhitelistMiddleware(allowedIPs))

		// Prometheus指标端点
		metricsGroup.GET("", metrics.MetricsHandler()) // 修改为空字符串，匹配/metrics路径

		// 健康检查端点
		metricsGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"version": "1.0.0",
			})
		})

		// 就绪检查端点
		metricsGroup.GET("/ready", func(c *gin.Context) {
			// 检查数据库连接
			container := c.MustGet("container").(*container.Container)
			db, err := container.DB().DB()
			if err != nil {
				c.JSON(503, gin.H{
					"status": "not_ready",
					"reason": "database_connection_error",
				})
				return
			}

			if err := db.Ping(); err != nil {
				c.JSON(503, gin.H{
					"status": "not_ready",
					"reason": "database_connection_error",
				})
				return
			}

			// 检查Redis连接
			redisClient := container.RedisClient()
			if redisClient != nil {
				if err := redisClient.Ping(c.Request.Context()); err != nil {
					c.JSON(503, gin.H{
						"status": "not_ready",
						"reason": "redis_connection_error",
					})
					return
				}
			}

			c.JSON(200, gin.H{
				"status": "ready",
			})
		})

		// 活跃性检查端点
		metricsGroup.GET("/live", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "alive",
			})
		})
	}
}
