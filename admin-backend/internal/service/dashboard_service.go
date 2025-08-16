package service

import (
	"context"
	"i18n-flow/internal/domain"
)

// DashboardService 仪表板服务实现
type DashboardService struct {
	projectRepo     domain.ProjectRepository
	languageRepo    domain.LanguageRepository
	translationRepo domain.TranslationRepository
}

// NewDashboardService 创建仪表板服务实例
func NewDashboardService(
	projectRepo domain.ProjectRepository,
	languageRepo domain.LanguageRepository,
	translationRepo domain.TranslationRepository,
) *DashboardService {
	return &DashboardService{
		projectRepo:     projectRepo,
		languageRepo:    languageRepo,
		translationRepo: translationRepo,
	}
}

// GetStats 获取仪表板统计信息
func (s *DashboardService) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	// 获取项目总数
	projects, totalProjects, err := s.projectRepo.GetAll(ctx, 1000000, 0) // 大数获取全部
	if err != nil {
		return nil, err
	}
	stats.TotalProjects = int(totalProjects)

	// 获取语言总数
	languages, err := s.languageRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalLanguages = len(languages)

	// 计算翻译总数和唯一键总数
	totalTranslations := 0
	uniqueKeys := make(map[string]bool)

	for _, project := range projects {
		// 获取每个项目的翻译矩阵
		matrix, err := s.translationRepo.GetMatrix(ctx, project.ID)
		if err != nil {
			continue // 跳过出错的项目
		}

		// 统计翻译数和唯一键
		for key, translations := range matrix {
			uniqueKeys[key] = true
			totalTranslations += len(translations)
		}
	}

	stats.TotalTranslations = totalTranslations
	stats.TotalKeys = len(uniqueKeys)

	return stats, nil
}
