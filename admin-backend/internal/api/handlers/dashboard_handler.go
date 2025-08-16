package handlers

import (
	"i18n-flow/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表板处理器
type DashboardHandler struct {
	dashboardService domain.DashboardService
}

// NewDashboardHandler 创建仪表板处理器
func NewDashboardHandler(dashboardService domain.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetStats 获取仪表板统计信息
// @Summary      获取仪表板统计信息
// @Description  获取项目、语言、翻译等统计信息
// @Tags         仪表板
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.DashboardStats
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /dashboard/stats [get]
func (h *DashboardHandler) GetStats(ctx *gin.Context) {
	stats, err := h.dashboardService.GetStats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
