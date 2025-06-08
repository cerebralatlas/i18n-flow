package main

import (
	"i18n-flow/config"
	"i18n-flow/controller"
	"i18n-flow/docs" // 导入自动生成的 docs 包
	"i18n-flow/middleware"
	"i18n-flow/model/db"
	"i18n-flow/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
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
	// 获取配置
	appConfig := config.GetConfig()

	// 初始化多日志系统
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	logConfig := utils.MultiLogConfig{
		Level:         appConfig.Log.Level,
		Format:        appConfig.Log.Format,
		Output:        appConfig.Log.Output,
		LogDir:        appConfig.Log.LogDir,
		DateFormat:    appConfig.Log.DateFormat,
		MaxSize:       appConfig.Log.MaxSize,
		MaxAge:        appConfig.Log.MaxAge,
		MaxBackups:    appConfig.Log.MaxBackups,
		Compress:      appConfig.Log.Compress,
		EnableConsole: appConfig.Log.EnableConsole,
		LogTypes: map[string]string{
			"access": "info",
			"error":  "error",
			"auth":   "info",
			"db":     "warn",
			"app":    appConfig.Log.Level,
		},
	}

	if err := utils.InitMultiLogger(logConfig); err != nil {
		panic("Failed to initialize multi-logger: " + err.Error())
	}

	// 确保程序退出时同步所有日志
	defer utils.SyncAll()

	utils.AppInfo("Application starting",
		zap.String("version", "1.0.0"),
		zap.String("environment", env),
		zap.String("log_level", appConfig.Log.Level),
		zap.String("log_dir", appConfig.Log.LogDir),
	)

	// 初始化数据库连接
	utils.AppInfo("Initializing database connection")
	db.InitDB()

	// 初始化Swagger文档
	docs.SwaggerInfo.BasePath = "/api"

	router := gin.Default()

	// 移除Gin默认的日志中间件
	router.Use(gin.Recovery())

	// 请求ID中间件
	router.Use(middleware.RequestIDMiddleware())

	// 请求日志中间件（跳过健康检查）
	router.Use(middleware.SkipLoggingMiddleware("/health"))
	router.Use(middleware.LoggingMiddleware())

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

	utils.AppInfo("Server starting",
		zap.String("address", ":8080"),
		zap.String("docs", "http://localhost:8080/swagger/index.html"),
	)

	if err := router.Run(":8080"); err != nil {
		utils.AppError("Failed to start server", zap.Error(err))
		os.Exit(1)
	}
}
