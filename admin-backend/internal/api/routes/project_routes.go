package routes

import "github.com/gin-gonic/gin"

// setupProjectRoutes 设置项目相关路由
func (r *Router) setupProjectRoutes(authRoutes *gin.RouterGroup) {
	projectRoutes := authRoutes.Group("/projects")
	{
		// 项目基本操作
		projectRoutes.POST("", r.ProjectHandler.Create)
		projectRoutes.GET("", r.ProjectHandler.GetAll)
		projectRoutes.GET("/accessible", r.ProjectHandler.GetAccessibleProjects)

		// 需要项目查看权限的操作
		projectViewRoutes := projectRoutes.Group("")
		projectViewRoutes.Use(r.middlewareFactory.RequireProjectViewer())
		{
			projectViewRoutes.GET("/detail/:id", r.ProjectHandler.GetByID)
			projectViewRoutes.GET("/:project_id/members", r.ProjectMemberHandler.GetProjectMembers)
			projectViewRoutes.GET("/:project_id/members/:user_id/permission", r.ProjectMemberHandler.CheckPermission)
		}

		// 需要项目编辑权限的操作
		projectEditRoutes := projectRoutes.Group("")
		projectEditRoutes.Use(r.middlewareFactory.RequireProjectEditor())
		{
			projectEditRoutes.PUT("/update/:id", r.ProjectHandler.Update)
		}

		// 需要项目所有者权限的操作
		projectOwnerRoutes := projectRoutes.Group("")
		projectOwnerRoutes.Use(r.middlewareFactory.RequireProjectOwner())
		{
			projectOwnerRoutes.DELETE("/delete/:id", r.ProjectHandler.Delete)
			projectOwnerRoutes.POST("/:project_id/members", r.ProjectMemberHandler.AddMember)
			projectOwnerRoutes.PUT("/:project_id/members/:user_id", r.ProjectMemberHandler.UpdateMemberRole)
			projectOwnerRoutes.DELETE("/:project_id/members/:user_id", r.ProjectMemberHandler.RemoveMember)
		}
	}
}
