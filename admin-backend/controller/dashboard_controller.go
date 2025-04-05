package controller

import (
	"i18n-flow/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardController 处理仪表板相关的请求
type DashboardController struct {
	dashboardService service.DashboardService
}

// ErrorResponse 表示错误响应的结构
type ErrorResponse struct {
	Error string `json:"error" example:"错误信息"`
}

// NewDashboardController 创建一个新的仪表板控制器
func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: service.NewDashboardService(),
	}
}

// GetDashboardStats 获取仪表板统计数据
// @Summary 获取仪表板统计数据
// @Description 获取系统的各种统计数据，包括项目总数、翻译总数、语言总数和用户总数
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} service.DashboardStats
// @Failure 500 {object} ErrorResponse
// @Router /dashboard/stats [get]
// @Security BearerAuth
func (dc *DashboardController) GetDashboardStats(c *gin.Context) {
	stats, err := dc.dashboardService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "获取统计数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
