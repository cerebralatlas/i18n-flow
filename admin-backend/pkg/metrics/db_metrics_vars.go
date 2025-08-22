package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// 数据库连接池指标
	dbConnectionsOpen = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_db_connections_open",
			Help: "Number of open connections in the database pool",
		},
	)

	dbConnectionsInUse = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_db_connections_in_use",
			Help: "Number of connections currently in use in the database pool",
		},
	)

	dbConnectionsIdle = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_db_connections_idle",
			Help: "Number of idle connections in the database pool",
		},
	)

	dbConnectionsWaitCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_db_connections_wait_count",
			Help: "Total number of connections waited for",
		},
	)

	dbConnectionsWaitDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_db_connections_wait_duration_ms",
			Help: "Total time waited for connections in milliseconds",
		},
	)

	dbConnectionsMaxIdleClosed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_db_connections_max_idle_closed",
			Help: "Total number of connections closed due to SetMaxIdleConns",
		},
	)

	dbConnectionsMaxLifetimeClosed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "i18n_flow_db_connections_max_lifetime_closed",
			Help: "Total number of connections closed due to SetConnMaxLifetime",
		},
	)

	// 查询类型计数器
	dbQueryTypeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "i18n_flow_db_query_type_total",
			Help: "Total number of database queries by type",
		},
		[]string{"type"},
	)

	// 查询错误计数器
	dbQueryErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "i18n_flow_db_query_errors_total",
			Help: "Total number of database query errors by type",
		},
		[]string{"type", "error"},
	)

	// 数据库指标是否已注册
	dbMetricsRegistered = false
	dbMetricsLock       = &sync.Mutex{}
)

// registerDBConnectionMetrics 注册数据库连接池指标
func registerDBConnectionMetrics() {
	dbMetricsLock.Lock()
	defer dbMetricsLock.Unlock()

	if dbMetricsRegistered {
		return
	}

	// 注册数据库连接池指标
	prometheus.MustRegister(
		dbConnectionsOpen,
		dbConnectionsInUse,
		dbConnectionsIdle,
		dbConnectionsWaitCount,
		dbConnectionsWaitDuration,
		dbConnectionsMaxIdleClosed,
		dbConnectionsMaxLifetimeClosed,
		dbQueryTypeCounter,
		dbQueryErrorCounter,
	)

	dbMetricsRegistered = true
}

// RecordDBQueryType 记录数据库查询类型
func RecordDBQueryType(queryType string) {
	dbQueryTypeCounter.WithLabelValues(queryType).Inc()
}

// RecordDBQueryError 记录数据库查询错误
func RecordDBQueryError(queryType, errorType string) {
	dbQueryErrorCounter.WithLabelValues(queryType, errorType).Inc()
}
