package metrics

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// DBMetricsCollector 数据库指标收集器
type DBMetricsCollector struct {
	interval time.Duration
	stopCh   chan struct{}
	db       *gorm.DB
}

// NewDBMetricsCollector 创建数据库指标收集器
func NewDBMetricsCollector(interval time.Duration, db *gorm.DB) *DBMetricsCollector {
	return &DBMetricsCollector{
		interval: interval,
		stopCh:   make(chan struct{}),
		db:       db,
	}
}

// Start 开始收集数据库指标
func (c *DBMetricsCollector) Start() {
	go func() {
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

// Stop 停止收集数据库指标
func (c *DBMetricsCollector) Stop() {
	close(c.stopCh)
}

// collectMetrics 收集数据库指标
func (c *DBMetricsCollector) collectMetrics() {
	// 收集数据库连接池统计信息
	c.collectConnectionStats()
}

// collectConnectionStats 收集数据库连接池统计信息
func (c *DBMetricsCollector) collectConnectionStats() {
	// 获取数据库连接池统计信息
	db, err := c.db.DB()
	if err != nil {
		return
	}

	// 注册数据库连接池指标（如果尚未注册）
	registerDBConnectionMetrics()

	// 设置连接池统计信息
	dbConnectionsOpen.Set(float64(db.Stats().OpenConnections))
	dbConnectionsInUse.Set(float64(db.Stats().InUse))
	dbConnectionsIdle.Set(float64(db.Stats().Idle))
	dbConnectionsWaitCount.Set(float64(db.Stats().WaitCount))
	dbConnectionsWaitDuration.Set(float64(db.Stats().WaitDuration.Milliseconds()))
	dbConnectionsMaxIdleClosed.Set(float64(db.Stats().MaxIdleClosed))
	dbConnectionsMaxLifetimeClosed.Set(float64(db.Stats().MaxLifetimeClosed))
}

// DBQueryHook GORM查询钩子，用于收集数据库查询指标
type DBQueryHook struct{}

// Before 在查询执行前调用
func (h *DBQueryHook) Before(ctx context.Context, stmt *gorm.Statement) (context.Context, error) {
	// 在上下文中存储开始时间
	return context.WithValue(ctx, "query_start_time", time.Now()), nil
}

// After 在查询执行后调用
func (h *DBQueryHook) After(ctx context.Context, stmt *gorm.Statement) {
	// 获取开始时间
	startTime, ok := ctx.Value("query_start_time").(time.Time)
	if !ok {
		return
	}

	// 计算查询耗时
	duration := time.Since(startTime)

	// 获取操作类型和实体名称
	operation := getOperationType(stmt)
	entity := getEntityName(stmt)

	// 记录数据库操作指标
	RecordDBOperation(operation, entity, duration)
}

// getOperationType 获取操作类型
func getOperationType(stmt *gorm.Statement) string {
	switch {
	case stmt.SQL.String() == "":
		return "unknown"
	case stmt.SQL.Len() >= 6 && stmt.SQL.String()[0:6] == "SELECT":
		return "select"
	case stmt.SQL.Len() >= 6 && stmt.SQL.String()[0:6] == "INSERT":
		return "insert"
	case stmt.SQL.Len() >= 6 && stmt.SQL.String()[0:6] == "UPDATE":
		return "update"
	case stmt.SQL.Len() >= 6 && stmt.SQL.String()[0:6] == "DELETE":
		return "delete"
	default:
		return "other"
	}
}

// getEntityName 获取实体名称
func getEntityName(stmt *gorm.Statement) string {
	if stmt.Schema != nil {
		return stmt.Schema.Table
	}
	return "unknown"
}
