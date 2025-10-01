package container

import (
	"i18n-flow/internal/config"
	"i18n-flow/internal/domain"
	"i18n-flow/internal/repository"
	"i18n-flow/internal/service"

	"gorm.io/gorm"
)

// Container 依赖注入容器
type Container struct {
	db     *gorm.DB
	config *config.Config

	// Repositories
	userRepo        domain.UserRepository
	projectRepo     domain.ProjectRepository
	languageRepo    domain.LanguageRepository
	translationRepo domain.TranslationRepository
	redisClient     *repository.RedisClient

	// Services
	authService        domain.AuthService
	userService        domain.UserService
	projectService     domain.ProjectService
	languageService    domain.LanguageService
	translationService domain.TranslationService
	dashboardService   domain.DashboardService
	cacheService       domain.CacheService
}

// NewContainer 创建新的容器实例
func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	return &Container{
		db:     db,
		config: cfg,
	}
}

// DB 获取数据库连接
func (c *Container) DB() *gorm.DB {
	return c.db
}

// Config 获取配置
func (c *Container) Config() *config.Config {
	return c.config
}

// UserRepository 获取用户仓储
func (c *Container) UserRepository() domain.UserRepository {
	if c.userRepo == nil {
		c.userRepo = c.createUserRepository()
	}
	return c.userRepo
}

// ProjectRepository 获取项目仓储
func (c *Container) ProjectRepository() domain.ProjectRepository {
	if c.projectRepo == nil {
		c.projectRepo = c.createProjectRepository()
	}
	return c.projectRepo
}

// LanguageRepository 获取语言仓储
func (c *Container) LanguageRepository() domain.LanguageRepository {
	if c.languageRepo == nil {
		c.languageRepo = c.createLanguageRepository()
	}
	return c.languageRepo
}

// TranslationRepository 获取翻译仓储
func (c *Container) TranslationRepository() domain.TranslationRepository {
	if c.translationRepo == nil {
		c.translationRepo = c.createTranslationRepository()
	}
	return c.translationRepo
}

// AuthService 获取认证服务
func (c *Container) AuthService() domain.AuthService {
	if c.authService == nil {
		c.authService = c.createAuthService()
	}
	return c.authService
}

// UserService 获取用户服务
func (c *Container) UserService() domain.UserService {
	if c.userService == nil {
		c.userService = c.createUserService()
	}
	return c.userService
}

// ProjectService 获取项目服务
func (c *Container) ProjectService() domain.ProjectService {
	if c.projectService == nil {
		c.projectService = c.createProjectService()
	}
	return c.projectService
}

// LanguageService 获取语言服务
func (c *Container) LanguageService() domain.LanguageService {
	if c.languageService == nil {
		c.languageService = c.createLanguageService()
	}
	return c.languageService
}

// TranslationService 获取翻译服务
func (c *Container) TranslationService() domain.TranslationService {
	if c.translationService == nil {
		c.translationService = c.createTranslationService()
	}
	return c.translationService
}

// DashboardService 获取仪表板服务
func (c *Container) DashboardService() domain.DashboardService {
	if c.dashboardService == nil {
		c.dashboardService = c.createDashboardService()
	}
	return c.dashboardService
}

// RedisClient 获取Redis客户端
func (c *Container) RedisClient() *repository.RedisClient {
	if c.redisClient == nil {
		c.redisClient = c.createRedisClient()
	}
	return c.redisClient
}

// CacheService 获取缓存服务
func (c *Container) CacheService() domain.CacheService {
	if c.cacheService == nil {
		c.cacheService = c.createCacheService()
	}
	return c.cacheService
}

// 创建仓储实例的私有方法

// createUserRepository 创建用户仓储
func (c *Container) createUserRepository() domain.UserRepository {
	return repository.NewUserRepository(c.db)
}

// createProjectRepository 创建项目仓储
func (c *Container) createProjectRepository() domain.ProjectRepository {
	return repository.NewProjectRepository(c.db)
}

// createLanguageRepository 创建语言仓储
func (c *Container) createLanguageRepository() domain.LanguageRepository {
	return repository.NewLanguageRepository(c.db)
}

// createTranslationRepository 创建翻译仓储
func (c *Container) createTranslationRepository() domain.TranslationRepository {
	return repository.NewTranslationRepository(c.db)
}

// createRedisClient 创建Redis客户端
func (c *Container) createRedisClient() *repository.RedisClient {
	return repository.NewRedisClient(&c.config.Redis)
}

// 创建服务实例的私有方法

// createAuthService 创建认证服务
func (c *Container) createAuthService() domain.AuthService {
	return service.NewAuthService(c.config.JWT)
}

// createUserService 创建用户服务
func (c *Container) createUserService() domain.UserService {
	baseService := service.NewUserService(c.UserRepository(), c.AuthService())
	// 如果有缓存服务，返回带缓存的版本
	if c.CacheService() != nil {
		return service.NewCachedUserService(baseService, c.CacheService())
	}
	return baseService
}

// createProjectService 创建项目服务
func (c *Container) createProjectService() domain.ProjectService {
	baseService := service.NewProjectService(c.ProjectRepository())
	// 如果有缓存服务，返回带缓存的版本
	if c.CacheService() != nil {
		return service.NewCachedProjectService(baseService, c.CacheService())
	}
	return baseService
}

// createLanguageService 创建语言服务
func (c *Container) createLanguageService() domain.LanguageService {
	baseService := service.NewLanguageService(c.LanguageRepository())
	// 如果有缓存服务，返回带缓存的版本
	if c.CacheService() != nil {
		return service.NewCachedLanguageService(baseService, c.CacheService())
	}
	return baseService
}

// createTranslationService 创建翻译服务
func (c *Container) createTranslationService() domain.TranslationService {
	baseService := service.NewTranslationService(
		c.TranslationRepository(),
		c.ProjectRepository(),
		c.LanguageRepository(),
	)
	// 如果有缓存服务，返回带缓存的版本
	if c.CacheService() != nil {
		return service.NewCachedTranslationService(baseService, c.CacheService())
	}
	return baseService
}

// createDashboardService 创建仪表板服务
func (c *Container) createDashboardService() domain.DashboardService {
	baseService := service.NewDashboardService(
		c.ProjectRepository(),
		c.LanguageRepository(),
		c.TranslationRepository(),
	)
	// 如果有缓存服务，返回带缓存的版本
	if c.CacheService() != nil {
		return service.NewCachedDashboardService(baseService, c.CacheService())
	}
	return baseService
}

// createCacheService 创建缓存服务
func (c *Container) createCacheService() domain.CacheService {
	return service.NewCacheService(c.RedisClient())
}
