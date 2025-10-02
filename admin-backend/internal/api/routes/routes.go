package routes

import (
	"i18n-flow/internal/api/handlers"
	"i18n-flow/internal/api/middleware"
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/container"
	internal_utils "i18n-flow/internal/utils"
	"i18n-flow/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// Router 路由器
type Router struct {
	container         *container.Container
	handlerFactory    *handlers.HandlerFactory
	middlewareFactory *middleware.MiddlewareFactory
}

// NewRouter 创建路由器
func NewRouter(container *container.Container) *Router {
	return &Router{
		container:      container,
		handlerFactory: handlers.NewHandlerFactory(container),
		middlewareFactory: middleware.NewMiddlewareFactory(
			container.AuthService(),
			container.UserService(),
			container.ProjectMemberService(),
		),
	}
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes(engine *gin.Engine, monitor *internal_utils.SimpleMonitor) {
	// 基本路由
	engine.GET("/", func(c *gin.Context) {
		response.Success(c, gin.H{"message": "Hello, World!"})
	})

	// 监控端点
	r.setupMonitoringRoutes(engine, monitor)

	// Swagger 文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := engine.Group("/api")
	{
		r.setupPublicRoutes(api)
		r.setupAuthenticatedRoutes(api)
		r.setupCLIRoutes(api)
	}
}

// setupPublicRoutes 设置公开路由
func (r *Router) setupPublicRoutes(rg *gin.RouterGroup) {
	userHandler := r.handlerFactory.UserHandler()

	// 登录路由组（应用登录限流中间件）
	loginRoutes := rg.Group("")
	loginRoutes.Use(middleware.TollboothLoginRateLimitMiddleware())
	{
		// 公开的认证路由（每秒5个请求，突发10个）
		loginRoutes.POST("/login", userHandler.Login)
		loginRoutes.POST("/refresh", userHandler.RefreshToken)
	}
}

// setupAuthenticatedRoutes 设置需要认证的路由
func (r *Router) setupAuthenticatedRoutes(rg *gin.RouterGroup) {
	// 应用JWT认证中间件和API限流中间件
	authRoutes := rg.Group("")
	authRoutes.Use(r.middlewareFactory.JWTAuthMiddleware())
	authRoutes.Use(middleware.TollboothAPIRateLimitMiddleware())

	// 获取处理器
	userHandler := r.handlerFactory.UserHandler()
	projectHandler := r.handlerFactory.ProjectHandler()
	languageHandler := r.handlerFactory.LanguageHandler()
	translationHandler := r.handlerFactory.TranslationHandler()
	dashboardHandler := r.handlerFactory.DashboardHandler()
	projectMemberHandler := r.handlerFactory.ProjectMemberHandler()

	// 用户相关路由
	userRoutes := authRoutes.Group("/user")
	{
		userRoutes.GET("/info", userHandler.GetUserInfo)
		userRoutes.POST("/change-password", userHandler.ChangePassword)
	}

	// 用户管理路由（管理员功能）
	usersRoutes := authRoutes.Group("/users")
	usersRoutes.Use(r.middlewareFactory.RequireAdminRole()) // 用户管理需要管理员权限
	{
		usersRoutes.POST("", userHandler.CreateUser)
		usersRoutes.GET("", userHandler.GetUsers)
		usersRoutes.GET("/:id", userHandler.GetUser)
		usersRoutes.PUT("/:id", userHandler.UpdateUser)
		usersRoutes.POST("/:id/reset-password", userHandler.ResetPassword)
		usersRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	// 用户项目关联路由（单独的路由组避免冲突）
	userProjectRoutes := authRoutes.Group("/user-projects")
	userProjectRoutes.Use(r.middlewareFactory.RequireAdminRole())
	{
		userProjectRoutes.GET("/:user_id", projectMemberHandler.GetUserProjects)
	}

	// 项目相关路由
	projectRoutes := authRoutes.Group("/projects")
	{
		// 项目基本操作
		projectRoutes.POST("", projectHandler.Create)
		projectRoutes.GET("", projectHandler.GetAll)
		projectRoutes.GET("/accessible", projectHandler.GetAccessibleProjects)

		// 需要项目查看权限的操作
		projectViewRoutes := projectRoutes.Group("")
		projectViewRoutes.Use(r.middlewareFactory.RequireProjectViewer())
		{
			projectViewRoutes.GET("/detail/:id", projectHandler.GetByID)
			projectViewRoutes.GET("/:project_id/members", projectMemberHandler.GetProjectMembers)
			projectViewRoutes.GET("/:project_id/members/:user_id/permission", projectMemberHandler.CheckPermission)
		}

		// 需要项目编辑权限的操作
		projectEditRoutes := projectRoutes.Group("")
		projectEditRoutes.Use(r.middlewareFactory.RequireProjectEditor())
		{
			projectEditRoutes.PUT("/update/:id", projectHandler.Update)
		}

		// 需要项目所有者权限的操作
		projectOwnerRoutes := projectRoutes.Group("")
		projectOwnerRoutes.Use(r.middlewareFactory.RequireProjectOwner())
		{
			projectOwnerRoutes.DELETE("/delete/:id", projectHandler.Delete)
			projectOwnerRoutes.POST("/:project_id/members", projectMemberHandler.AddMember)
			projectOwnerRoutes.PUT("/:project_id/members/:user_id", projectMemberHandler.UpdateMemberRole)
			projectOwnerRoutes.DELETE("/:project_id/members/:user_id", projectMemberHandler.RemoveMember)
		}
	}

	// 语言相关路由
	languageRoutes := authRoutes.Group("/languages")
	{
		languageRoutes.GET("", languageHandler.GetAll) // 所有用户都可以查看语言列表

		// 语言管理需要管理员权限
		languageAdminRoutes := languageRoutes.Group("")
		languageAdminRoutes.Use(r.middlewareFactory.RequireAdminRole())
		{
			languageAdminRoutes.POST("", languageHandler.Create)
			languageAdminRoutes.PUT("/:id", languageHandler.Update)
			languageAdminRoutes.DELETE("/:id", languageHandler.Delete)
		}
	}

	// 翻译相关路由
	translationRoutes := authRoutes.Group("/translations")
	{
		// 需要项目查看权限的操作
		translationViewRoutes := translationRoutes.Group("")
		translationViewRoutes.Use(r.middlewareFactory.RequireProjectViewer())
		{
			translationViewRoutes.GET("/by-project/:project_id", translationHandler.GetByProjectID)
			translationViewRoutes.GET("/matrix/by-project/:project_id", translationHandler.GetMatrix)
			translationViewRoutes.GET("/:id", translationHandler.GetByID)
		}

		// 需要项目编辑权限的操作
		translationEditRoutes := translationRoutes.Group("")
		translationEditRoutes.Use(r.middlewareFactory.RequireProjectEditor())
		{
			translationEditRoutes.POST("", translationHandler.Create)
			translationEditRoutes.PUT("/:id", translationHandler.Update)
			translationEditRoutes.DELETE("/:id", translationHandler.Delete)
		}
	}

	// 批量操作路由组（应用批量操作限流中间件和项目编辑权限）
	batchRoutes := authRoutes.Group("/translations")
	batchRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
	batchRoutes.Use(r.middlewareFactory.RequireProjectEditor())
	{
		batchRoutes.POST("/batch", translationHandler.CreateBatch)
		batchRoutes.POST("/batch-delete", translationHandler.DeleteBatch)
	}

	// 导出导入路由（应用批量操作限流中间件和项目权限）
	exportRoutes := authRoutes.Group("/exports")
	exportRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
	exportRoutes.Use(r.middlewareFactory.RequireProjectViewer()) // 导出只需要查看权限
	{
		exportRoutes.GET("/project/:project_id", translationHandler.Export)
	}

	importRoutes := authRoutes.Group("/imports")
	importRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
	importRoutes.Use(r.middlewareFactory.RequireProjectEditor()) // 导入需要编辑权限
	{
		importRoutes.POST("/project/:project_id", translationHandler.Import)
	}

	// 仪表板相关路由
	dashboardRoutes := authRoutes.Group("/dashboard")
	{
		dashboardRoutes.GET("/stats", dashboardHandler.GetStats)
	}
}

// setupCLIRoutes 设置CLI相关路由
func (r *Router) setupCLIRoutes(rg *gin.RouterGroup) {
	// CLI路由使用API Key认证和API限流
	cliRoutes := rg.Group("/cli")
	cliRoutes.Use(r.middlewareFactory.APIKeyAuthMiddleware())
	cliRoutes.Use(middleware.APIRateLimitMiddleware())

	// 获取CLI处理器
	cliHandler := r.handlerFactory.CLIHandler()

	{
		// CLI身份验证
		cliRoutes.GET("/auth", cliHandler.Auth)

		// 获取翻译数据
		cliRoutes.GET("/translations", cliHandler.GetTranslations)
	}

	// 推送翻译键（批量操作，应用批量操作限流）
	batchCliRoutes := rg.Group("/cli")
	batchCliRoutes.Use(r.middlewareFactory.APIKeyAuthMiddleware())
	batchCliRoutes.Use(middleware.BatchOperationRateLimitMiddleware())
	{
		batchCliRoutes.POST("/keys", cliHandler.PushKeys)
	}
}

// setupMonitoringRoutes 设置监控路由
func (r *Router) setupMonitoringRoutes(engine *gin.Engine, monitor *internal_utils.SimpleMonitor) {
	// 健康检查端点（替换原有的简单健康检查）
	engine.GET("/health", monitor.HealthCheck)

	// 基础统计端点
	engine.GET("/stats", monitor.SimpleStats)

	// 详细统计端点
	engine.GET("/stats/detailed", monitor.DetailedStats)

	utils.AppInfo("Monitoring endpoints configured",
		zap.String("health_check", "GET /health"),
		zap.String("basic_stats", "GET /stats"),
		zap.String("detailed_stats", "GET /stats/detailed"),
	)
}
