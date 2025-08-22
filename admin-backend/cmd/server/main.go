package main

import (
	"context"
	"i18n-flow/docs" // 导入自动生成的 docs 包
	"i18n-flow/internal/api/middleware"
	"i18n-flow/internal/api/routes"
	"i18n-flow/internal/config"
	"i18n-flow/internal/container"
	"i18n-flow/internal/repository"
	"i18n-flow/pkg/metrics"
	"i18n-flow/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志系统
	if err := initLogger(cfg); err != nil {
		log.Fatalf("初始化日志系统失败: %v", err)
	}
	defer utils.SyncAll()

	// 获取环境变量
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	utils.AppInfo("Application starting",
		zap.String("version", "1.0.0"),
		zap.String("environment", env),
		zap.String("log_level", cfg.Log.Level),
		zap.String("log_dir", cfg.Log.LogDir),
	)

	// 初始化数据库连接
	utils.AppInfo("Initializing database connection")
	db, err := repository.InitDB(cfg)
	if err != nil {
		utils.AppError("Failed to initialize database", zap.Error(err))
		os.Exit(1)
	}

	// 创建依赖注入容器
	container := container.NewContainer(db, cfg)

	// 初始化Redis连接
	utils.AppInfo("Initializing Redis connection")
	redisClient := container.RedisClient()
	if err := redisClient.Ping(context.Background()); err != nil {
		utils.AppWarn("Failed to connect to Redis, cache will be disabled", zap.Error(err))
	} else {
		utils.AppInfo("Redis connection established successfully")
	}

	// 初始化Swagger文档
	docs.SwaggerInfo.BasePath = "/api"

	// 创建Gin引擎
	router := gin.Default()

	// 移除Gin默认的日志中间件
	router.Use(gin.Recovery())

	// 应用全局中间件
	setupMiddleware(router)

	// 设置路由
	routeManager := routes.NewRouter(container)
	routeManager.SetupRoutes(router)

	// 设置监控路由
	routes.SetupMetricsRoutes(router, container)

	// 初始化监控系统
	metricsEnabled := false
	if cfg.Metrics.Enabled {
		metricsEnabled = true
	}

	if metricsEnabled {
		utils.AppInfo("Initializing metrics system")
		metrics.InitMetrics(container)
	}

	// 启动服务器
	utils.AppInfo("Server starting",
		zap.String("address", ":8080"),
		zap.String("docs", "http://localhost:8080/swagger/index.html"),
	)

	if err := router.Run(":8080"); err != nil {
		utils.AppError("Failed to start server", zap.Error(err))
		os.Exit(1)
	}
}

// initLogger 初始化日志系统
func initLogger(cfg *config.Config) error {
	logConfig := utils.MultiLogConfig{
		Level:         cfg.Log.Level,
		Format:        cfg.Log.Format,
		Output:        cfg.Log.Output,
		LogDir:        cfg.Log.LogDir,
		DateFormat:    cfg.Log.DateFormat,
		MaxSize:       cfg.Log.MaxSize,
		MaxAge:        cfg.Log.MaxAge,
		MaxBackups:    cfg.Log.MaxBackups,
		Compress:      cfg.Log.Compress,
		EnableConsole: cfg.Log.EnableConsole,
		LogTypes: map[string]string{
			"access": "info",
			"error":  "error",
			"auth":   "info",
			"db":     "warn",
			"app":    cfg.Log.Level,
		},
	}

	return utils.InitMultiLogger(logConfig)
}

// setupMiddleware 设置全局中间件
func setupMiddleware(router *gin.Engine) {
	// 安全HTTP头中间件（最先设置，确保所有响应都包含安全头）
	router.Use(middleware.SecurityHeadersMiddleware())

	// 全局限流中间件（每秒100个请求，突发200个）
	router.Use(middleware.GlobalRateLimitMiddleware())

	// 请求ID中间件
	router.Use(middleware.RequestIDMiddleware())

	// 请求日志中间件（跳过健康检查）
	router.Use(middleware.SkipLoggingMiddleware("/health"))
	router.Use(middleware.LoggingMiddleware())

	// 全局错误处理中间件
	router.Use(middleware.ErrorHandlerMiddleware())

	// 请求大小限制中间件 (32MB)
	router.Use(middleware.RequestSizeLimitMiddleware(32 << 20))

	// 请求验证中间件
	router.Use(middleware.RequestValidationMiddleware())

	// 分页参数验证中间件
	router.Use(middleware.PaginationValidationMiddleware())

	// 允许跨域请求
	router.Use(middleware.CORSMiddleware())

	// 404处理器
	router.NoRoute(middleware.NotFoundHandler())
}
