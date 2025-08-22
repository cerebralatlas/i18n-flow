package metrics

import (
	"i18n-flow/internal/container"
	"time"
)

// 全局指标收集器
var (
	systemMetricsCollector   *SystemMetricsCollector
	dbMetricsCollector       *DBMetricsCollector
	redisMetricsCollector    *RedisMetricsCollector
	businessMetricsCollector *BusinessMetricsCollector
	alertManager             *AlertManager
)

// InitMetrics 初始化监控系统
func InitMetrics(c *container.Container) {
	// 初始化Prometheus指标
	Init()

	// 创建系统指标收集器（每30秒收集一次）
	systemMetricsCollector = NewSystemMetricsCollector(30 * time.Second)
	systemMetricsCollector.Start()

	// 创建数据库指标收集器（每30秒收集一次）
	dbMetricsCollector = NewDBMetricsCollector(30*time.Second, c.DB())
	dbMetricsCollector.Start()

	// 创建Redis指标收集器（每30秒收集一次）
	if c.RedisClient() != nil {
		redisMetricsCollector = NewRedisMetricsCollector(30*time.Second, c.RedisClient())
		redisMetricsCollector.Start()
	}

	// 创建业务指标收集器（每分钟收集一次）
	businessMetricsCollector = NewBusinessMetricsCollector(60*time.Second, c.Repository())
	businessMetricsCollector.Start()

	// 创建告警管理器
	initAlertManager()
}

// initAlertManager 初始化告警管理器
func initAlertManager() {
	// 创建告警配置
	alertConfig := AlertConfig{
		CheckInterval:             60 * time.Second,
		SendRecoveryNotifications: true,
		SuppressWindow:            30 * time.Minute,
	}

	// 创建告警管理器
	alertManager = NewAlertManager(alertConfig)

	// 添加告警阈值
	addDefaultAlertThresholds(alertManager)

	// 添加告警通知器
	// alertManager.AddNotifier(&WebhookNotifier{URL: "http://alert-webhook.example.com"})

	// 启动告警管理器
	alertManager.Start()
}

// addDefaultAlertThresholds 添加默认告警阈值
func addDefaultAlertThresholds(am *AlertManager) {
	// 添加高CPU使用率告警
	am.AddThreshold("high_cpu_usage", AlertThreshold{
		MetricName:   "i18n_flow_cpu_usage_percent",
		MetricLabels: map[string]string{},
		Threshold:    90,
		Operator:     ">",
		Severity:     "warning",
		Description:  "CPU使用率超过90%",
	})

	// 添加高内存使用率告警
	am.AddThreshold("high_memory_usage", AlertThreshold{
		MetricName:   "i18n_flow_memory_usage_bytes",
		MetricLabels: map[string]string{},
		Threshold:    1073741824, // 1GB
		Operator:     ">",
		Severity:     "warning",
		Description:  "内存使用量超过1GB",
	})

	// 添加高HTTP错误率告警
	am.AddThreshold("high_http_error_rate", AlertThreshold{
		MetricName:   "i18n_flow_http_requests_total",
		MetricLabels: map[string]string{"status": "500"},
		Threshold:    10,
		Operator:     ">",
		Severity:     "error",
		Description:  "HTTP 500错误数超过10个",
	})

	// 添加数据库连接池告警
	am.AddThreshold("db_connection_pool_exhaustion", AlertThreshold{
		MetricName:   "i18n_flow_db_connections_in_use",
		MetricLabels: map[string]string{},
		Threshold:    90, // 假设连接池上限为100
		Operator:     ">",
		Severity:     "critical",
		Description:  "数据库连接池即将耗尽",
	})

	// 添加Redis连接告警
	am.AddThreshold("redis_connection_failure", AlertThreshold{
		MetricName:   "i18n_flow_redis_errors_total",
		MetricLabels: map[string]string{"operation": "ping", "error": "connection refused"},
		Threshold:    1,
		Operator:     ">=",
		Severity:     "critical",
		Description:  "Redis连接失败",
	})
}

// ShutdownMetrics 关闭监控系统
func ShutdownMetrics() {
	// 停止系统指标收集器
	if systemMetricsCollector != nil {
		systemMetricsCollector.Stop()
	}

	// 停止数据库指标收集器
	if dbMetricsCollector != nil {
		dbMetricsCollector.Stop()
	}

	// 停止Redis指标收集器
	if redisMetricsCollector != nil {
		redisMetricsCollector.Stop()
	}

	// 停止业务指标收集器
	if businessMetricsCollector != nil {
		businessMetricsCollector.Stop()
	}

	// 停止告警管理器
	if alertManager != nil {
		alertManager.Stop()
	}
}
