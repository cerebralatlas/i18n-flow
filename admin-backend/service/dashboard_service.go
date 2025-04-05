package service

import (
	"i18n-flow/model"
	"i18n-flow/model/db"
)

// DashboardService 提供仪表板相关的服务
type DashboardService struct{}

// NewDashboardService 创建一个新的仪表板服务
func NewDashboardService() DashboardService {
	return DashboardService{}
}

// DashboardStats 包含仪表板需要的统计数据
type DashboardStats struct {
	ProjectCount     int64 `json:"project_count"`
	TranslationCount int64 `json:"translation_count"`
	LanguageCount    int64 `json:"language_count"`
	UserCount        int64 `json:"user_count"`
}

// GetDashboardStats 获取仪表板的统计数据
func (ds DashboardService) GetDashboardStats() (DashboardStats, error) {
	database := db.DB
	var stats DashboardStats

	// 获取项目总数
	var projectCount int64
	if err := database.Model(&model.Project{}).Count(&projectCount).Error; err != nil {
		return stats, err
	}
	stats.ProjectCount = projectCount

	// 获取翻译键的总数（而不是所有翻译条目的总数）
	// 使用正确的列名 key_name 而不是 key
	var keyCount int64
	if err := database.Model(&model.Translation{}).
		Distinct("key_name").
		Count(&keyCount).Error; err != nil {
		return stats, err
	}
	stats.TranslationCount = keyCount

	// 获取语言总数
	var languageCount int64
	if err := database.Model(&model.Language{}).Count(&languageCount).Error; err != nil {
		return stats, err
	}
	stats.LanguageCount = languageCount

	// 获取用户总数
	var userCount int64
	if err := database.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return stats, err
	}
	stats.UserCount = userCount

	return stats, nil
}
