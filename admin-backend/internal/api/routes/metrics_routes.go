package routes

import (
	"i18n-flow/internal/api/middleware"
	"i18n-flow/internal/container"
	"i18n-flow/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// SetupMetricsRoutes 设置监控相关路由
func SetupMetricsRoutes(router *gin.Engine, c *container.Container) {
	// 监控指标路由组
	metricsGroup := router.Group("/metrics")
	{
		// 添加IP白名单中间件，限制只有特定IP可以访问监控指标
		metricsGroup.Use(middleware.IPWhitelistMiddleware(c.Config().Metrics.AllowedIPs))

		// Prometheus指标端点
		metricsGroup.GET("/", metrics.MetricsHandler())

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
			db, err := c.MustGet("container").(*container.Container).DB().DB()
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
			redisClient := c.MustGet("container").(*container.Container).RedisClient()
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
