package container

import (
	"i18n-flow/internal/domain"
	"i18n-flow/internal/repository"
	"i18n-flow/internal/service"
)

// createUserRepository 创建用户仓储实例
func (c *Container) createUserRepository() domain.UserRepository {
	return repository.NewUserRepository(c.db)
}

// createProjectRepository 创建项目仓储实例
func (c *Container) createProjectRepository() domain.ProjectRepository {
	return repository.NewProjectRepository(c.db)
}

// createLanguageRepository 创建语言仓储实例
func (c *Container) createLanguageRepository() domain.LanguageRepository {
	return repository.NewLanguageRepository(c.db)
}

// createTranslationRepository 创建翻译仓储实例
func (c *Container) createTranslationRepository() domain.TranslationRepository {
	return repository.NewTranslationRepository(c.db)
}

// RedisClient 获取Redis客户端
func (c *Container) RedisClient() *repository.RedisClient {
	if c.redisClient == nil {
		c.redisClient = repository.NewRedisClient(&c.config.Redis)
	}
	return c.redisClient
}

// CacheService 获取缓存服务
func (c *Container) CacheService() domain.CacheService {
	if c.cacheService == nil {
		c.cacheService = service.NewCacheService(c.RedisClient())
	}
	return c.cacheService
}

// createAuthService 创建认证服务实例
func (c *Container) createAuthService() domain.AuthService {
	return service.NewAuthService(c.config.JWT)
}

// createUserService 创建用户服务实例
func (c *Container) createUserService() domain.UserService {
	return service.NewUserService(
		c.UserRepository(),
		c.AuthService(),
	)
}

// createProjectService 创建项目服务实例
func (c *Container) createProjectService() domain.ProjectService {
	return service.NewProjectService(c.ProjectRepository())
}

// createLanguageService 创建语言服务实例
func (c *Container) createLanguageService() domain.LanguageService {
	return service.NewLanguageService(c.LanguageRepository())
}

// createTranslationService 创建翻译服务实例
func (c *Container) createTranslationService() domain.TranslationService {
	baseService := service.NewTranslationService(
		c.TranslationRepository(),
		c.ProjectRepository(),
		c.LanguageRepository(),
	)

	// 使用缓存装饰器包装基础服务
	return service.NewCachedTranslationService(baseService, c.CacheService())
}

// createDashboardService 创建仪表板服务实例
func (c *Container) createDashboardService() domain.DashboardService {
	baseService := service.NewDashboardService(
		c.ProjectRepository(),
		c.LanguageRepository(),
		c.TranslationRepository(),
	)

	// 使用缓存装饰器包装基础服务
	return service.NewCachedDashboardService(baseService, c.CacheService())
}
