package container

import (
	"i18n-flow/internal/config"
	"i18n-flow/internal/domain"
	"i18n-flow/internal/repository"

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
