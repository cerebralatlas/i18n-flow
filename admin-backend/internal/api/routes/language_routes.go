package routes

import "github.com/gin-gonic/gin"

// setupLanguageRoutes 设置语言相关路由
func (r *Router) setupLanguageRoutes(authRoutes *gin.RouterGroup) {
	languageRoutes := authRoutes.Group("/languages")
	{
		languageRoutes.GET("", r.LanguageHandler.GetAll) // 所有用户都可以查看语言列表

		// 语言管理需要管理员权限
		languageAdminRoutes := languageRoutes.Group("")
		languageAdminRoutes.Use(r.middlewareFactory.RequireAdminRole())
		{
			languageAdminRoutes.POST("", r.LanguageHandler.Create)
			languageAdminRoutes.PUT("/:id", r.LanguageHandler.Update)
			languageAdminRoutes.DELETE("/:id", r.LanguageHandler.Delete)
		}
	}
}
