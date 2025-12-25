package routes

import (
	"i18n-flow/internal/api/handlers"
	"i18n-flow/internal/api/middleware"
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"
	internal_utils "i18n-flow/internal/utils"
	"i18n-flow/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Router 路由器
type Router struct {
	UserHandler          *handlers.UserHandler
	ProjectHandler       *handlers.ProjectHandler
	LanguageHandler      *handlers.LanguageHandler
	TranslationHandler   *handlers.TranslationHandler
	DashboardHandler     *handlers.DashboardHandler
	ProjectMemberHandler *handlers.ProjectMemberHandler
	CLIHandler           *handlers.CLIHandler
	middlewareFactory    *middleware.MiddlewareFactory
}

// RouterDeps 定义 Router 的依赖（用于 fx.In）
type RouterDeps struct {
	fx.In
	UserHandler          *handlers.UserHandler
	ProjectHandler       *handlers.ProjectHandler
	LanguageHandler      *handlers.LanguageHandler
	TranslationHandler   *handlers.TranslationHandler
	DashboardHandler     *handlers.DashboardHandler
	ProjectMemberHandler *handlers.ProjectMemberHandler
	CLIHandler           *handlers.CLIHandler
	AuthService          domain.AuthService
	UserService          domain.UserService
	ProjectMemberService domain.ProjectMemberService
}

// NewRouter 创建路由器
func NewRouter(deps RouterDeps) *Router {
	return &Router{
		UserHandler:          deps.UserHandler,
		ProjectHandler:       deps.ProjectHandler,
		LanguageHandler:      deps.LanguageHandler,
		TranslationHandler:   deps.TranslationHandler,
		DashboardHandler:     deps.DashboardHandler,
		ProjectMemberHandler: deps.ProjectMemberHandler,
		CLIHandler:           deps.CLIHandler,
		middlewareFactory: middleware.NewMiddlewareFactory(
			deps.AuthService,
			deps.UserService,
			deps.ProjectMemberService,
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
	// 登录路由组（应用登录限流中间件）
	loginRoutes := rg.Group("")
	loginRoutes.Use(middleware.TollboothLoginRateLimitMiddleware())
	{
		// 公开的认证路由（每秒5个请求，突发10个）
		loginRoutes.POST("/login", r.UserHandler.Login)
		loginRoutes.POST("/refresh", r.UserHandler.RefreshToken)
	}
}

// setupAuthenticatedRoutes 设置需要认证的路由
func (r *Router) setupAuthenticatedRoutes(rg *gin.RouterGroup) {
	// 应用JWT认证中间件和API限流中间件
	authRoutes := rg.Group("")
	authRoutes.Use(r.middlewareFactory.JWTAuthMiddleware())
	authRoutes.Use(middleware.TollboothAPIRateLimitMiddleware())

	// 用户相关路由
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

	// 项目相关路由
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

	// 语言相关路由
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

	// 翻译相关路由
	translationRoutes := authRoutes.Group("/translations")
	{
		// 需要项目查看权限的操作
		translationViewRoutes := translationRoutes.Group("")
		translationViewRoutes.Use(r.middlewareFactory.RequireProjectViewer())
		{
			translationViewRoutes.GET("/by-project/:project_id", r.TranslationHandler.GetByProjectID)
			translationViewRoutes.GET("/matrix/by-project/:project_id", r.TranslationHandler.GetMatrix)
			translationViewRoutes.GET("/:id", r.TranslationHandler.GetByID)
		}

		// 需要项目编辑权限的操作
		translationEditRoutes := translationRoutes.Group("")
		translationEditRoutes.Use(r.middlewareFactory.RequireProjectEditor())
		{
			translationEditRoutes.POST("", r.TranslationHandler.Create)
			translationEditRoutes.PUT("/:id", r.TranslationHandler.Update)
			translationEditRoutes.DELETE("/:id", r.TranslationHandler.Delete)
		}
	}

	// 批量操作路由组（应用批量操作限流中间件和项目编辑权限）
	batchRoutes := authRoutes.Group("/translations")
	batchRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
	batchRoutes.Use(r.middlewareFactory.RequireProjectEditor())
	{
		batchRoutes.POST("/batch", r.TranslationHandler.CreateBatch)
		batchRoutes.POST("/batch-delete", r.TranslationHandler.DeleteBatch)
	}

	// 导出导入路由（应用批量操作限流中间件和项目权限）
	exportRoutes := authRoutes.Group("/exports")
	exportRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
	exportRoutes.Use(r.middlewareFactory.RequireProjectViewer()) // 导出只需要查看权限
	{
		exportRoutes.GET("/project/:project_id", r.TranslationHandler.Export)
	}

	importRoutes := authRoutes.Group("/imports")
	importRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
	importRoutes.Use(r.middlewareFactory.RequireProjectEditor()) // 导入需要编辑权限
	{
		importRoutes.POST("/project/:project_id", r.TranslationHandler.Import)
	}

	// 仪表板相关路由
	dashboardRoutes := authRoutes.Group("/dashboard")
	{
		dashboardRoutes.GET("/stats", r.DashboardHandler.GetStats)
	}
}

// setupCLIRoutes 设置CLI相关路由
func (r *Router) setupCLIRoutes(rg *gin.RouterGroup) {
	// CLI路由使用API Key认证和API限流
	cliRoutes := rg.Group("/cli")
	cliRoutes.Use(r.middlewareFactory.APIKeyAuthMiddleware())
	cliRoutes.Use(middleware.TollboothAPIRateLimitMiddleware())
	{
		// CLI身份验证
		cliRoutes.GET("/auth", r.CLIHandler.Auth)

		// 获取翻译数据
		cliRoutes.GET("/translations", r.CLIHandler.GetTranslations)
	}

	// 推送翻译键（批量操作，应用批量操作限流）
	batchCliRoutes := rg.Group("/cli")
	batchCliRoutes.Use(r.middlewareFactory.APIKeyAuthMiddleware())
	batchCliRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
	{
		batchCliRoutes.POST("/keys", r.CLIHandler.PushKeys)
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

// RouterModule 定义路由模块
var RouterModule = fx.Module("router",
	fx.Provide(NewRouter),
)
