package metrics

import (
	"context"
	"i18n-flow/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisHook Redis钩子，用于收集Redis操作指标
type RedisHook struct{}

// BeforeProcess 在命令执行前调用
func (h *RedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	// 在上下文中存储开始时间
	return context.WithValue(ctx, "redis_start_time", time.Now()), nil
}

// AfterProcess 在命令执行后调用
func (h *RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	// 获取开始时间
	startTime, ok := ctx.Value("redis_start_time").(time.Time)
	if !ok {
		return nil
	}

	// 计算操作耗时
	duration := time.Since(startTime)

	// 获取操作名称
	operation := cmd.Name()

	// 记录Redis操作指标
	RecordRedisOperation(operation, duration)

	// 如果命令执行出错，记录错误
	if err := cmd.Err(); err != nil && err != redis.Nil {
		RecordRedisError(operation, err.Error())
	}

	return nil
}

// BeforeProcessPipeline 在管道命令执行前调用
func (h *RedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	// 在上下文中存储开始时间
	return context.WithValue(ctx, "redis_pipeline_start_time", time.Now()), nil
}

// AfterProcessPipeline 在管道命令执行后调用
func (h *RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	// 获取开始时间
	startTime, ok := ctx.Value("redis_pipeline_start_time").(time.Time)
	if !ok {
		return nil
	}

	// 计算操作耗时
	duration := time.Since(startTime)

	// 记录Redis管道操作指标
	RecordRedisOperation("pipeline", duration)

	// 记录管道中的各个命令
	for _, cmd := range cmds {
		operation := cmd.Name()

		// 记录Redis操作指标（不计算单个命令的耗时）
		redisOperationCounter.WithLabelValues(operation).Inc()

		// 如果命令执行出错，记录错误
		if err := cmd.Err(); err != nil && err != redis.Nil {
			RecordRedisError(operation, err.Error())
		}
	}

	return nil
}

// 实现redis.Hook接口所需的方法
func (h *RedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// RedisMetricsCollector Redis指标收集器
type RedisMetricsCollector struct {
	interval    time.Duration
	stopCh      chan struct{}
	redisClient *repository.RedisClient
}

// NewRedisMetricsCollector 创建Redis指标收集器
func NewRedisMetricsCollector(interval time.Duration, redisClient *repository.RedisClient) *RedisMetricsCollector {
	return &RedisMetricsCollector{
		interval:    interval,
		stopCh:      make(chan struct{}),
		redisClient: redisClient,
	}
}

// Start 开始收集Redis指标
func (c *RedisMetricsCollector) Start() {
	// 由于我们使用的是自定义的RedisClient封装，不能直接添加钩子
	// 这里我们只能通过定期收集Redis信息来获取指标

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

// Stop 停止收集Redis指标
func (c *RedisMetricsCollector) Stop() {
	close(c.stopCh)
}

// collectMetrics 收集Redis指标
func (c *RedisMetricsCollector) collectMetrics() {
	ctx := context.Background()

	// 检查Redis连接
	if err := c.redisClient.Ping(ctx); err != nil {
		RecordRedisError("ping", err.Error())
		return
	}

	// 由于我们的RedisClient是封装的，无法直接获取底层的Redis统计信息
	// 这里只能通过一些简单的操作来评估Redis状态

	// 记录Redis操作
	startTime := time.Now()
	_, _ = c.redisClient.Get(ctx, "metrics_health_check")
	RecordRedisOperation("get", time.Since(startTime))

	// 设置一个键来测试Redis写入性能
	startTime = time.Now()
	_ = c.redisClient.Set(ctx, "metrics_health_check", time.Now().String(), time.Minute)
	RecordRedisOperation("set", time.Since(startTime))
}
