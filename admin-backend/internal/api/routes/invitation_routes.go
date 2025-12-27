package routes

import "github.com/gin-gonic/gin"

// setupInvitationRoutes 设置邀请相关路由
func (r *Router) setupInvitationRoutes(authRoutes *gin.RouterGroup) {
	// 邀请管理路由（管理员功能）
	invitationRoutes := authRoutes.Group("/invitations")
	invitationRoutes.Use(r.middlewareFactory.RequireAdminRole()) // 邀请管理需要管理员权限
	{
		invitationRoutes.POST("", r.InvitationHandler.CreateInvitation)
		invitationRoutes.GET("", r.InvitationHandler.GetInvitations)
		invitationRoutes.GET("/:code", r.InvitationHandler.GetInvitation)
		invitationRoutes.DELETE("/:code", r.InvitationHandler.RevokeInvitation)
	}
}

// setupPublicInvitationRoutes 设置公开的邀请相关路由
func (r *Router) setupPublicInvitationRoutes(rg *gin.RouterGroup) {
	// 公开的邀请验证和注册路由（不需要认证）
	publicInvitationRoutes := rg.Group("/invitations")
	{
		publicInvitationRoutes.GET("/:code/validate", r.InvitationHandler.ValidateInvitation)
	}
}

// setupPublicRegisterRoutes 设置公开的注册路由
func (r *Router) setupPublicRegisterRoutes(rg *gin.RouterGroup) {
	// 公开的注册路由（不需要认证）
	rg.POST("/register", r.InvitationHandler.RegisterWithInvitation)
}
