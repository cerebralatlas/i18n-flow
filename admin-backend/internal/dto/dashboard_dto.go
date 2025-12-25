package dto

// DashboardStats 仪表板统计
type DashboardStats struct {
	TotalProjects     int `json:"total_projects"`
	TotalLanguages    int `json:"total_languages"`
	TotalTranslations int `json:"total_translations"`
	TotalKeys         int `json:"total_keys"`
}
