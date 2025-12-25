package di

import (
	"i18n-flow/internal/api/handlers"
	"i18n-flow/internal/config"

	"go.uber.org/fx"
)

// AppModule 定义主模块
var AppModule = fx.Module("app",
	// 配置提供者 - 从外部传入的 cfg 获取
	fx.Provide(func(cfg *config.Config) *config.Config {
		return cfg
	}),

	// 数据库和缓存
	fx.Provide(NewDB),
	fx.Provide(NewRedisClient),

	// 缓存服务
	fx.Provide(NewCacheService),

	// Repositories
	fx.Provide(NewUserRepository),
	fx.Provide(NewProjectRepository),
	fx.Provide(NewLanguageRepository),
	fx.Provide(NewTranslationRepository),
	fx.Provide(NewProjectMemberRepository),

	// Auth Service (无缓存)
	fx.Provide(NewAuthService),

	// Services (带缓存装饰器)
	fx.Provide(NewUserService),
	fx.Provide(NewProjectService),
	fx.Provide(NewLanguageService),
	fx.Provide(NewTranslationService),
	fx.Provide(NewDashboardService),
	fx.Provide(NewProjectMemberService),

	// Handlers
	fx.Provide(handlers.NewUserHandler),
	fx.Provide(handlers.NewProjectHandler),
	fx.Provide(handlers.NewLanguageHandler),
	fx.Provide(handlers.NewTranslationHandler),
	fx.Provide(handlers.NewProjectMemberHandler),
	fx.Provide(handlers.NewCLIHandler),
	fx.Provide(handlers.NewDashboardHandler),
)
