package routes

import (
	"i18n-flow/internal/api/handlers"
	"i18n-flow/internal/api/middleware"
	"i18n-flow/internal/container"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		container:         container,
		handlerFactory:    handlers.NewHandlerFactory(container),
		middlewareFactory: middleware.NewMiddlewareFactory(container.AuthService()),
	}
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes(engine *gin.Engine) {
	// 基本路由
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// 健康检查端点
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "i18n-flow",
			"version": "1.0.0",
		})
	})

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

	// 公开的认证路由
	rg.POST("/login", userHandler.Login)
	rg.POST("/refresh", userHandler.RefreshToken)
}

// setupAuthenticatedRoutes 设置需要认证的路由
func (r *Router) setupAuthenticatedRoutes(rg *gin.RouterGroup) {
	// 应用JWT认证中间件
	authRoutes := rg.Group("")
	authRoutes.Use(r.middlewareFactory.JWTAuthMiddleware())

	// 获取处理器
	userHandler := r.handlerFactory.UserHandler()
	projectHandler := r.handlerFactory.ProjectHandler()
	languageHandler := r.handlerFactory.LanguageHandler()
	translationHandler := r.handlerFactory.TranslationHandler()
	dashboardHandler := r.handlerFactory.DashboardHandler()

	// 用户相关路由
	userRoutes := authRoutes.Group("/user")
	{
		userRoutes.GET("/info", userHandler.GetUserInfo)
	}

	// 项目相关路由
	projectRoutes := authRoutes.Group("/projects")
	{
		projectRoutes.POST("", projectHandler.Create)
		projectRoutes.GET("", projectHandler.GetAll)
		projectRoutes.GET("/detail/:id", projectHandler.GetByID)
		projectRoutes.PUT("/update/:id", projectHandler.Update)
		projectRoutes.DELETE("/delete/:id", projectHandler.Delete)
	}

	// 语言相关路由
	languageRoutes := authRoutes.Group("/languages")
	{
		languageRoutes.GET("", languageHandler.GetAll)
		languageRoutes.POST("", languageHandler.Create)
		languageRoutes.PUT("/:id", languageHandler.Update)
		languageRoutes.DELETE("/:id", languageHandler.Delete)
	}

	// 翻译相关路由
	translationRoutes := authRoutes.Group("/translations")
	{
		translationRoutes.POST("", translationHandler.Create)
		translationRoutes.POST("/batch", translationHandler.CreateBatch)
		translationRoutes.GET("/by-project/:project_id", translationHandler.GetByProjectID)
		translationRoutes.GET("/matrix/by-project/:project_id", translationHandler.GetMatrix)
		translationRoutes.GET("/:id", translationHandler.GetByID)
		translationRoutes.PUT("/:id", translationHandler.Update)
		translationRoutes.DELETE("/:id", translationHandler.Delete)
		translationRoutes.POST("/batch-delete", translationHandler.DeleteBatch)
	}

	// 导出导入路由
	exportRoutes := authRoutes.Group("/exports")
	{
		exportRoutes.GET("/project/:project_id", translationHandler.Export)
	}

	// 仪表板相关路由
	dashboardRoutes := authRoutes.Group("/dashboard")
	{
		dashboardRoutes.GET("/stats", dashboardHandler.GetStats)
	}
}

// setupCLIRoutes 设置CLI相关路由
func (r *Router) setupCLIRoutes(rg *gin.RouterGroup) {
	// CLI路由需要API Key认证，这里暂时跳过
	// 可以根据需要实现API Key中间件
	cliRoutes := rg.Group("/cli")
	{
		// 这里可以添加CLI相关的路由
		_ = cliRoutes
	}
}
