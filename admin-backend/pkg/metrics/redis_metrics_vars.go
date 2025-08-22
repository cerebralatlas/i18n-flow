package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	// Redis错误计数器
	redisErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "i18n_flow_redis_errors_total",
			Help: "Total number of Redis errors by operation and error type",
		},
		[]string{"operation", "error"},
	)

	// Redis连接池指标
	redisPoolStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "i18n_flow_redis_pool_stats",
			Help: "Redis connection pool statistics",
		},
		[]string{"stat"},
	)

	// Redis内存使用指标
	redisMemoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_redis_memory_used_bytes",
			Help: "Redis memory usage in bytes",
		},
	)

	// Redis键数量指标
	redisKeyCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_redis_keys_total",
			Help: "Total number of keys in Redis",
		},
	)

	// Redis过期键数量指标
	redisExpiredKeyCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_redis_expired_keys_total",
			Help: "Total number of expired keys in Redis",
		},
	)

	// Redis驱逐键数量指标
	redisEvictedKeyCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_redis_evicted_keys_total",
			Help: "Total number of evicted keys in Redis",
		},
	)

	// Redis连接数指标
	redisConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_redis_connections",
			Help: "Number of client connections to Redis",
		},
	)

	// Redis命令处理数指标
	redisCommandsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "i18n_flow_redis_commands_processed_total",
			Help: "Total number of commands processed by Redis",
		},
	)

	// Redis命令错误数指标
	redisCommandErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "i18n_flow_redis_command_errors_total",
			Help: "Total number of command errors in Redis",
		},
	)

	// Redis指标是否已注册
	redisMetricsRegistered = false
	redisMetricsLock       = &sync.Mutex{}
)

// registerRedisMetrics 注册Redis指标
func registerRedisMetrics() {
	redisMetricsLock.Lock()
	defer redisMetricsLock.Unlock()

	if redisMetricsRegistered {
		return
	}

	// 注册Redis指标
	prometheus.MustRegister(
		redisErrorCounter,
		redisPoolStats,
		redisMemoryUsage,
		redisKeyCount,
		redisExpiredKeyCount,
		redisEvictedKeyCount,
		redisConnections,
		redisCommandsProcessed,
		redisCommandErrors,
	)

	redisMetricsRegistered = true
}

// RecordRedisError 记录Redis错误
func RecordRedisError(operation, errorType string) {
	registerRedisMetrics()
	redisErrorCounter.WithLabelValues(operation, errorType).Inc()
}

// SetRedisPoolStats 设置Redis连接池统计信息
func SetRedisPoolStats(stat string, value float64) {
	registerRedisMetrics()
	redisPoolStats.WithLabelValues(stat).Set(value)
}

// SetRedisMemoryUsage 设置Redis内存使用量
func SetRedisMemoryUsage(bytes float64) {
	registerRedisMetrics()
	redisMemoryUsage.Set(bytes)
}

// SetRedisKeyCount 设置Redis键数量
func SetRedisKeyCount(count float64) {
	registerRedisMetrics()
	redisKeyCount.Set(count)
}

// SetRedisExpiredKeyCount 设置Redis过期键数量
func SetRedisExpiredKeyCount(count float64) {
	registerRedisMetrics()
	redisExpiredKeyCount.Set(count)
}

// SetRedisEvictedKeyCount 设置Redis驱逐键数量
func SetRedisEvictedKeyCount(count float64) {
	registerRedisMetrics()
	redisEvictedKeyCount.Set(count)
}

// SetRedisConnections 设置Redis连接数
func SetRedisConnections(count float64) {
	registerRedisMetrics()
	redisConnections.Set(count)
}

// IncRedisCommandsProcessed 增加Redis命令处理数
func IncRedisCommandsProcessed() {
	registerRedisMetrics()
	redisCommandsProcessed.Inc()
}

// IncRedisCommandErrors 增加Redis命令错误数
func IncRedisCommandErrors() {
	registerRedisMetrics()
	redisCommandErrors.Inc()
}
