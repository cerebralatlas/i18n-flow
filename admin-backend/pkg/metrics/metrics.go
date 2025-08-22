package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 请求计数器
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "i18n_flow_http_requests_total",
			Help: "Total number of HTTP requests by status code, method and path",
		},
		[]string{"status", "method", "path"},
	)

	// 请求延迟直方图
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "i18n_flow_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// 当前活跃请求数
	activeRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "i18n_flow_http_active_requests",
			Help: "Number of active HTTP requests",
		},
		[]string{"method", "path"},
	)

	// 数据库操作计数器
	dbOperationCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "i18n_flow_db_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "entity"},
	)

	// 数据库操作延迟直方图
	dbOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "i18n_flow_db_operation_duration_seconds",
			Help:    "Database operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "entity"},
	)

	// Redis操作计数器
	redisOperationCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "i18n_flow_redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation"},
	)

	// Redis操作延迟直方图
	redisOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "i18n_flow_redis_operation_duration_seconds",
			Help:    "Redis operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// 系统内存使用量
	memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)

	// 系统CPU使用率
	cpuUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_cpu_usage_percent",
			Help: "Current CPU usage in percent",
		},
	)

	// 业务指标：翻译项目数量
	projectsCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_projects_count",
			Help: "Total number of translation projects",
		},
	)

	// 业务指标：翻译条目数量
	translationsCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "i18n_flow_translations_count",
			Help: "Total number of translations by project",
		},
		[]string{"project_id"},
	)

	// 业务指标：API调用次数
	apiCallCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "i18n_flow_api_calls_total",
			Help: "Total number of API calls by endpoint and client",
		},
		[]string{"endpoint", "client"},
	)

	// 业务指标：翻译完成率
	translationCompletionRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "i18n_flow_translation_completion_rate",
			Help: "Translation completion rate by project and language",
		},
		[]string{"project_id", "language"},
	)
)

// Init 初始化Prometheus指标
func Init() {
	// 注册所有指标
	prometheus.MustRegister(
		requestCounter,
		requestDuration,
		activeRequests,
		dbOperationCounter,
		dbOperationDuration,
		redisOperationCounter,
		redisOperationDuration,
		memoryUsage,
		cpuUsage,
		projectsCount,
		translationsCount,
		apiCallCounter,
		translationCompletionRate,
	)
}

// PrometheusMiddleware Gin中间件，用于收集HTTP请求指标
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method

		// 增加活跃请求计数
		activeRequests.WithLabelValues(method, path).Inc()

		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 减少活跃请求计数
		activeRequests.WithLabelValues(method, path).Dec()

		// 记录请求延迟
		duration := time.Since(startTime).Seconds()
		requestDuration.WithLabelValues(method, path).Observe(duration)

		// 记录请求状态码
		status := c.Writer.Status()
		requestCounter.WithLabelValues(string(rune(status)), method, path).Inc()
	}
}

// MetricsHandler 返回Prometheus指标的HTTP处理器
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// RecordDBOperation 记录数据库操作
func RecordDBOperation(operation, entity string, duration time.Duration) {
	dbOperationCounter.WithLabelValues(operation, entity).Inc()
	dbOperationDuration.WithLabelValues(operation, entity).Observe(duration.Seconds())
}

// RecordRedisOperation 记录Redis操作
func RecordRedisOperation(operation string, duration time.Duration) {
	redisOperationCounter.WithLabelValues(operation).Inc()
	redisOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// SetMemoryUsage 设置内存使用量
func SetMemoryUsage(bytes float64) {
	memoryUsage.Set(bytes)
}

// SetCPUUsage 设置CPU使用率
func SetCPUUsage(percent float64) {
	cpuUsage.Set(percent)
}

// SetProjectsCount 设置项目数量
func SetProjectsCount(count float64) {
	projectsCount.Set(count)
}

// SetTranslationsCount 设置翻译条目数量
func SetTranslationsCount(projectID string, count float64) {
	translationsCount.WithLabelValues(projectID).Set(count)
}

// RecordAPICall 记录API调用
func RecordAPICall(endpoint, client string) {
	apiCallCounter.WithLabelValues(endpoint, client).Inc()
}

// SetTranslationCompletionRate 设置翻译完成率
func SetTranslationCompletionRate(projectID, language string, rate float64) {
	translationCompletionRate.WithLabelValues(projectID, language).Set(rate)
}
