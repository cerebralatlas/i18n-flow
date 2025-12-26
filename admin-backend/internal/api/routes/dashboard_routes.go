package routes

import "github.com/gin-gonic/gin"

// setupDashboardRoutes 设置仪表板相关路由
func (r *Router) setupDashboardRoutes(authRoutes *gin.RouterGroup) {
	dashboardRoutes := authRoutes.Group("/dashboard")
	{
		dashboardRoutes.GET("/stats", r.DashboardHandler.GetStats)
	}
}
