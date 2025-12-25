package di

import (
	"fmt"

	"i18n-flow/internal/config"
	"i18n-flow/internal/domain"
	"i18n-flow/internal/repository"
	"i18n-flow/internal/service"

	"gorm.io/gorm"
)

// NewDB 提供数据库连接
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := repository.InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}
	return db, nil
}

// NewRedisClient 提供 Redis 客户端
func NewRedisClient(cfg *config.Config) *repository.RedisClient {
	return repository.NewRedisClient(&cfg.Redis)
}

// NewCacheService 提供缓存服务
func NewCacheService(client *repository.RedisClient) domain.CacheService {
	return service.NewCacheService(client)
}

// NewUserRepository 提供用户仓储
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return repository.NewUserRepository(db)
}

// NewProjectRepository 提供项目仓储
func NewProjectRepository(db *gorm.DB) domain.ProjectRepository {
	return repository.NewProjectRepository(db)
}

// NewLanguageRepository 提供语言仓储
func NewLanguageRepository(db *gorm.DB) domain.LanguageRepository {
	return repository.NewLanguageRepository(db)
}

// NewTranslationRepository 提供翻译仓储
func NewTranslationRepository(db *gorm.DB) domain.TranslationRepository {
	return repository.NewTranslationRepository(db)
}

// NewProjectMemberRepository 提供项目成员仓储
func NewProjectMemberRepository(db *gorm.DB) domain.ProjectMemberRepository {
	return repository.NewProjectMemberRepository(db)
}

// NewAuthService 提供认证服务
func NewAuthService(cfg *config.Config) domain.AuthService {
	return service.NewAuthService(cfg.JWT)
}

// NewUserService 提供用户服务 (带缓存装饰器)
func NewUserService(
	repo domain.UserRepository,
	auth domain.AuthService,
	cache domain.CacheService,
) domain.UserService {
	base := service.NewUserService(repo, auth)
	if cache != nil {
		return service.NewCachedUserService(base, cache)
	}
	return base
}

// NewProjectService 提供项目服务 (带缓存装饰器)
func NewProjectService(
	projectRepo domain.ProjectRepository,
	userRepo domain.UserRepository,
	memberRepo domain.ProjectMemberRepository,
	cache domain.CacheService,
) domain.ProjectService {
	base := service.NewProjectService(projectRepo, userRepo, memberRepo)
	if cache != nil {
		return service.NewCachedProjectService(base, cache)
	}
	return base
}

// NewLanguageService 提供语言服务 (带缓存装饰器)
func NewLanguageService(
	repo domain.LanguageRepository,
	cache domain.CacheService,
) domain.LanguageService {
	base := service.NewLanguageService(repo)
	if cache != nil {
		return service.NewCachedLanguageService(base, cache)
	}
	return base
}

// NewTranslationService 提供翻译服务 (带缓存装饰器)
func NewTranslationService(
	translationRepo domain.TranslationRepository,
	projectRepo domain.ProjectRepository,
	languageRepo domain.LanguageRepository,
	cache domain.CacheService,
) domain.TranslationService {
	base := service.NewTranslationService(translationRepo, projectRepo, languageRepo)
	if cache != nil {
		return service.NewCachedTranslationService(base, cache)
	}
	return base
}

// NewDashboardService 提供仪表板服务 (带缓存装饰器)
func NewDashboardService(
	projectRepo domain.ProjectRepository,
	languageRepo domain.LanguageRepository,
	translationRepo domain.TranslationRepository,
	cache domain.CacheService,
) domain.DashboardService {
	base := service.NewDashboardService(projectRepo, languageRepo, translationRepo)
	if cache != nil {
		return service.NewCachedDashboardService(base, cache)
	}
	return base
}

// NewProjectMemberService 提供项目成员服务
func NewProjectMemberService(
	memberRepo domain.ProjectMemberRepository,
	userRepo domain.UserRepository,
	projectRepo domain.ProjectRepository,
) domain.ProjectMemberService {
	return service.NewProjectMemberService(memberRepo, userRepo, projectRepo)
}
