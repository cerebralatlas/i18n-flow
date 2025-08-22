package metrics

import (
	"i18n-flow/internal/container"
	"time"
)

// BusinessMetricsCollector 业务指标收集器
type BusinessMetricsCollector struct {
	interval   time.Duration
	stopCh     chan struct{}
	repository *container.Repository
}

// NewBusinessMetricsCollector 创建业务指标收集器
func NewBusinessMetricsCollector(interval time.Duration, repo *container.Repository) *BusinessMetricsCollector {
	return &BusinessMetricsCollector{
		interval:   interval,
		stopCh:     make(chan struct{}),
		repository: repo,
	}
}

// Start 开始收集业务指标
func (c *BusinessMetricsCollector) Start() {
	go func() {
		// 立即收集一次
		c.collectMetrics()

		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.collectMetrics()
			case <-c.stopCh:
				return
			}
		}
	}()
}

// Stop 停止收集业务指标
func (c *BusinessMetricsCollector) Stop() {
	close(c.stopCh)
}

// collectMetrics 收集业务指标
func (c *BusinessMetricsCollector) collectMetrics() {
	// 设置一些模拟的业务指标数据
	SetProjectsCount(10)
	SetTranslationsCount("project-1", 100)
	SetTranslationsCount("project-2", 200)
	SetTranslationCompletionRate("project-1", "en", 80)
	SetTranslationCompletionRate("project-1", "zh", 60)
	SetTranslationCompletionRate("project-2", "en", 90)
	SetTranslationCompletionRate("project-2", "zh", 70)

	// 记录一些模拟的API调用
	RecordAPICall("/api/projects", "web")
	RecordAPICall("/api/translations", "cli")
}
