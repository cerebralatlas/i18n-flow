package main

import (
	"i18n-flow/controller"
	"i18n-flow/docs" // 导入自动生成的 docs 包
	"i18n-flow/middleware"
	"i18n-flow/model/db"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           i18n-flow API
// @version         1.0
// @description     i18n-flow 是一个用于管理多语言翻译的系统。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}
func main() {
	// 初始化数据库连接
	db.InitDB()

	// 初始化Swagger文档
	docs.SwaggerInfo.BasePath = "/api"

	router := gin.Default()

	// 全局错误处理中间件
	router.Use(middleware.ErrorHandlerMiddleware())

	// 请求验证中间件
	router.Use(middleware.RequestValidationMiddleware())

	// 允许跨域请求
	router.Use(middleware.CORSMiddleware())

	// 404处理器
	router.NoRoute(middleware.NotFoundHandler())

	// 基本路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "i18n-flow",
			"version": "1.0.0",
		})
	})

	// Swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 用户相关路由
	userController := controller.NewUserController()
	projectController := controller.NewProjectController()
	translationController := controller.NewTranslationController()
	languageController := controller.NewLanguageController()
	dashboardController := controller.NewDashboardController()
	cliController := controller.NewCLIController()

	// 公开路由
	router.POST("/api/login", userController.Login)
	router.POST("/api/refresh", userController.RefreshToken)

	// 需要鉴权的路由
	authRoutes := router.Group("/api")
	authRoutes.Use(middleware.JWTAuthMiddleware())
	{
		// 用户相关
		authRoutes.GET("/user/info", userController.GetUserInfo)

		// 项目相关
		authRoutes.POST("/projects", projectController.CreateProject)
		authRoutes.GET("/projects", projectController.GetProjects)
		authRoutes.GET("/projects/detail/:id", projectController.GetProjectByID)
		authRoutes.PUT("/projects/update/:id", projectController.UpdateProject)
		authRoutes.DELETE("/projects/delete/:id", projectController.DeleteProject)

		// 翻译相关
		authRoutes.POST("/translations", translationController.CreateTranslation)
		authRoutes.POST("/translations/batch", translationController.BatchCreateTranslations)
		authRoutes.GET("/translations/by-project/:project_id", translationController.GetTranslationsByProject)
		authRoutes.GET("/translations/matrix/by-project/:project_id", translationController.GetTranslationMatrix)
		authRoutes.GET("/translations/:id", translationController.GetTranslationByID)
		authRoutes.PUT("/translations/:id", translationController.UpdateTranslation)
		authRoutes.DELETE("/translations/:id", translationController.DeleteTranslation)
		authRoutes.POST("/translations/batch-delete", translationController.BatchDeleteTranslations)
		authRoutes.GET("/exports/project/:project_id", translationController.ExportTranslations)
		authRoutes.POST("/imports/project/:project_id", translationController.ImportTranslations)

		// 语言相关
		authRoutes.GET("/languages", languageController.GetLanguages)
		authRoutes.POST("/languages", languageController.CreateLanguage)
		authRoutes.PUT("/languages/:id", languageController.UpdateLanguage)
		authRoutes.DELETE("/languages/:id", languageController.DeleteLanguage)

		// 仪表板相关
		authRoutes.GET("/dashboard/stats", dashboardController.GetDashboardStats)
	}

	// CLI 相关
	cliRoutes := router.Group("/api/cli")
	cliRoutes.Use(middleware.APIKeyAuthMiddleware())
	{
		cliRoutes.GET("/auth", cliController.CheckAPIKey)
		cliRoutes.GET("/translations", cliController.GetAllTranslations)
		cliRoutes.POST("/keys", cliController.PushKeys)
	}

	router.Run(":8080")
}
