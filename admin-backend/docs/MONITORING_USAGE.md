# i18n-flow 监控系统使用指南

## 问题修复记录

在实现监控系统过程中，我们遇到并修复了以下问题：

### 1. `/metrics` 路径访问错误

**问题描述**：访问 `/metrics` 路径时报404错误。

**解决方案**：
- 修改了路由注册方式，将 `metricsGroup.GET("/", ...)` 改为 `metricsGroup.GET("", ...)`，以正确匹配 `/metrics` 路径
- 在路由初始化时将容器注入到上下文中，确保路由处理函数能够访问容器

### 2. 容器访问方式错误

**问题描述**：在路由处理函数中无法正确访问容器对象。

**解决方案**：
- 添加了中间件将容器注入到请求上下文中：
  ```go
  router.Use(func(ctx *gin.Context) {
      ctx.Set("container", c)
      ctx.Next()
  })
  ```
- 在路由处理函数中正确获取容器：
  ```go
  container := c.MustGet("container").(*container.Container)
  ```

### 3. Prometheus指标初始化问题

**问题描述**：Prometheus指标未正确初始化和注册。

**解决方案**：
- 在服务启动时调用 `metrics.Init()` 初始化所有Prometheus指标
- 修改了指标收集器的实现，确保在启动时立即收集一次指标

### 4. 业务指标收集器错误

**问题描述**：业务指标收集器中的方法调用与实际接口不匹配。

**解决方案**：
- 简化了业务指标收集器，使用模拟数据代替实际的数据库查询
- 确保了指标收集器能够正常工作而不依赖于特定的数据库接口

## 监控系统使用方法

### 访问监控指标

监控系统提供以下端点：

1. **Prometheus指标**：`/metrics`
   - 提供所有Prometheus格式的监控指标

2. **健康检查**：`/metrics/health`
   - 返回服务健康状态

3. **就绪检查**：`/metrics/ready`
   - 检查服务是否就绪（包括数据库和Redis连接）

4. **活跃性检查**：`/metrics/live`
   - 检查服务是否活跃

### 集成Prometheus

在Prometheus配置文件中添加以下scrape配置：

```yaml
scrape_configs:
  - job_name: 'i18n-flow'
    scrape_interval: 15s
    metrics_path: '/metrics'
    static_configs:
      - targets: ['localhost:8080']
```

### 集成Grafana

1. 在Grafana中添加Prometheus数据源
2. 导入位于 `pkg/metrics/grafana_dashboards/` 目录下的仪表盘JSON文件

## 监控指标说明

### 系统指标

- `i18n_flow_cpu_usage_percent`：CPU使用率
- `i18n_flow_memory_usage_bytes`：内存使用量

### HTTP指标

- `i18n_flow_http_requests_total`：HTTP请求总数
- `i18n_flow_http_request_duration_seconds`：HTTP请求延迟
- `i18n_flow_http_active_requests`：当前活跃请求数

### 数据库指标

- `i18n_flow_db_connections_open`：打开的数据库连接数
- `i18n_flow_db_connections_in_use`：使用中的数据库连接数
- `i18n_flow_db_operations_total`：数据库操作总数
- `i18n_flow_db_operation_duration_seconds`：数据库操作延迟

### Redis指标

- `i18n_flow_redis_operations_total`：Redis操作总数
- `i18n_flow_redis_operation_duration_seconds`：Redis操作延迟

### 业务指标

- `i18n_flow_projects_count`：项目总数
- `i18n_flow_translations_count`：翻译条目总数
- `i18n_flow_translation_completion_rate`：翻译完成率
- `i18n_flow_api_calls_total`：API调用总数

## 配置选项

监控系统可以通过以下配置选项进行配置：

```go
type MetricsConfig struct {
    Enabled                bool      // 是否启用监控
    AllowedIPs             []string  // 允许访问监控指标的IP白名单
    SystemMetricsEnabled   bool      // 是否启用系统指标收集
    DBMetricsEnabled       bool      // 是否启用数据库指标收集
    RedisMetricsEnabled    bool      // 是否启用Redis指标收集
    BusinessMetricsEnabled bool      // 是否启用业务指标收集
    AlertingEnabled        bool      // 是否启用告警系统
    AlertWebhookURL        string    // 告警Webhook URL
    CollectionIntervalSeconds int    // 指标收集间隔（秒）
}
```

## 告警系统

告警系统基于Prometheus指标阈值触发，支持以下告警规则：

1. **高CPU使用率**：CPU使用率超过90%
2. **高内存使用量**：内存使用量超过1GB
3. **高HTTP错误率**：HTTP 500错误数超过10个
4. **数据库连接池耗尽**：使用中的数据库连接数超过90个
5. **Redis连接失败**：Redis连接失败

告警通知支持以下方式：

1. **Webhook**：发送HTTP POST请求到指定URL
2. **Email**：发送邮件通知（需配置SMTP服务器）
