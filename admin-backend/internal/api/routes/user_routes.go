package routes

import "github.com/gin-gonic/gin"

// setupUserRoutes 设置用户相关路由
func (r *Router) setupUserRoutes(authRoutes *gin.RouterGroup) {
	// 当前用户路由
	userRoutes := authRoutes.Group("/user")
	{
		userRoutes.GET("/info", r.UserHandler.GetUserInfo)
		userRoutes.POST("/change-password", r.UserHandler.ChangePassword)
	}

	// 用户管理路由（管理员功能）
	usersRoutes := authRoutes.Group("/users")
	usersRoutes.Use(r.middlewareFactory.RequireAdminRole()) // 用户管理需要管理员权限
	{
		usersRoutes.POST("", r.UserHandler.CreateUser)
		usersRoutes.GET("", r.UserHandler.GetUsers)
		usersRoutes.GET("/:id", r.UserHandler.GetUser)
		usersRoutes.PUT("/:id", r.UserHandler.UpdateUser)
		usersRoutes.POST("/:id/reset-password", r.UserHandler.ResetPassword)
		usersRoutes.DELETE("/:id", r.UserHandler.DeleteUser)
	}

	// 用户项目关联路由（单独的路由组避免冲突）
	userProjectRoutes := authRoutes.Group("/user-projects")
	userProjectRoutes.Use(r.middlewareFactory.RequireAdminRole())
	{
		userProjectRoutes.GET("/:user_id", r.ProjectMemberHandler.GetUserProjects)
	}
}
