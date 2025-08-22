# i18n-flow 监控系统文档

本文档描述了 i18n-flow 后端服务的监控系统设计和使用方法。

## 概述

i18n-flow 监控系统基于 Prometheus 和 Grafana 构建，提供了全面的系统和业务指标监控，以及灵活的告警配置。监控系统包括以下几个主要组件：

1. **Prometheus 指标收集**：收集系统、数据库、Redis 和业务指标
2. **Grafana 仪表盘**：可视化监控数据
3. **告警系统**：基于阈值的告警机制
4. **健康检查 API**：提供服务健康状态检查

## 配置监控系统

### 环境变量配置

可以通过以下环境变量配置监控系统：

```
# 监控系统基本配置
METRICS_ENABLED=true                      # 是否启用监控系统
METRICS_COLLECTION_INTERVAL_SECONDS=30    # 指标收集间隔（秒）

# 指标访问控制
METRICS_ALLOWED_IPS=127.0.0.1,::1         # 允许访问监控指标的IP白名单

# 指标类型配置
METRICS_SYSTEM_ENABLED=true               # 是否启用系统指标
METRICS_DB_ENABLED=true                   # 是否启用数据库指标
METRICS_REDIS_ENABLED=true                # 是否启用Redis指标
METRICS_BUSINESS_ENABLED=true             # 是否启用业务指标

# 告警系统配置
METRICS_ALERTING_ENABLED=false            # 是否启用告警系统
METRICS_ALERT_WEBHOOK_URL=                # 告警Webhook URL

# 告警邮件配置
METRICS_ALERT_EMAIL_ENABLED=false         # 是否启用邮件告警
METRICS_ALERT_SMTP_SERVER=smtp.example.com # SMTP服务器
METRICS_ALERT_SMTP_PORT=587               # SMTP端口
METRICS_ALERT_SMTP_USER=user@example.com  # SMTP用户名
METRICS_ALERT_SMTP_PASSWORD=password      # SMTP密码
METRICS_ALERT_FROM_EMAIL=alerts@example.com # 发件人邮箱
METRICS_ALERT_TO_EMAILS=admin@example.com # 收件人邮箱（多个用逗号分隔）
```

### 配置文件

也可以在配置文件中设置监控系统参数：

```yaml
metrics:
  enabled: true
  allowed_ips:
    - 127.0.0.1
    - ::1
  system_metrics_enabled: true
  db_metrics_enabled: true
  redis_metrics_enabled: true
  business_metrics_enabled: true
  alerting_enabled: false
  alert_webhook_url: ""
  alert_email:
    enabled: false
    smtp_server: smtp.example.com
    smtp_port: 587
    smtp_user: user@example.com
    smtp_password: password
    from_email: alerts@example.com
    to_emails:
      - admin@example.com
  collection_interval_seconds: 30
```

## 监控指标

### 系统指标

- `i18n_flow_cpu_usage_percent`: CPU 使用率
- `i18n_flow_memory_usage_bytes`: 内存使用量
- `i18n_flow_http_requests_total`: HTTP 请求总数
- `i18n_flow_http_request_duration_seconds`: HTTP 请求延迟
- `i18n_flow_http_active_requests`: 当前活跃请求数

### 数据库指标

- `i18n_flow_db_connections_open`: 打开的数据库连接数
- `i18n_flow_db_connections_in_use`: 使用中的数据库连接数
- `i18n_flow_db_connections_idle`: 空闲的数据库连接数
- `i18n_flow_db_connections_wait_count`: 等待连接的次数
- `i18n_flow_db_connections_wait_duration_ms`: 等待连接的总时间
- `i18n_flow_db_operations_total`: 数据库操作总数
- `i18n_flow_db_operation_duration_seconds`: 数据库操作延迟

### Redis 指标

- `i18n_flow_redis_operations_total`: Redis 操作总数
- `i18n_flow_redis_operation_duration_seconds`: Redis 操作延迟
- `i18n_flow_redis_errors_total`: Redis 错误总数
- `i18n_flow_redis_memory_used_bytes`: Redis 内存使用量
- `i18n_flow_redis_keys_total`: Redis 键总数

### 业务指标

- `i18n_flow_projects_count`: 项目总数
- `i18n_flow_translations_count`: 翻译条目总数
- `i18n_flow_translation_completion_rate`: 翻译完成率
- `i18n_flow_api_calls_total`: API 调用总数

## Grafana 仪表盘

系统提供了三个预配置的 Grafana 仪表盘：

1. **系统仪表盘**：显示系统级指标，如 CPU、内存使用率和 HTTP 请求统计
2. **数据库仪表盘**：显示数据库连接池和查询性能指标
3. **业务仪表盘**：显示业务相关指标，如项目数量、翻译条目数量和完成率

### 导入仪表盘

1. 登录 Grafana
2. 点击左侧菜单的 "+" 按钮，然后选择 "Import"
3. 上传仪表盘 JSON 文件（位于 `pkg/metrics/grafana_dashboards/` 目录下）或粘贴 JSON 内容
4. 选择 Prometheus 数据源
5. 点击 "Import" 完成导入

## 告警系统

告警系统基于 Prometheus 指标阈值触发，支持以下告警通知方式：

- **Webhook**：发送 HTTP POST 请求到指定 URL
- **Email**：发送邮件通知

### 默认告警规则

系统预配置了以下告警规则：

1. **高 CPU 使用率**：CPU 使用率超过 90%
2. **高内存使用量**：内存使用量超过 1GB
3. **高 HTTP 错误率**：HTTP 500 错误数超过 10 个
4. **数据库连接池耗尽**：使用中的数据库连接数超过 90 个
5. **Redis 连接失败**：Redis 连接失败

### 自定义告警规则

可以通过修改 `pkg/metrics/init.go` 文件中的 `addDefaultAlertThresholds` 函数来自定义告警规则。

## 健康检查 API

系统提供了以下健康检查 API：

- **GET /metrics/health**：返回服务健康状态
- **GET /metrics/ready**：返回服务就绪状态（检查数据库和 Redis 连接）
- **GET /metrics/live**：返回服务活跃状态

## Prometheus 集成

要将 i18n-flow 监控系统与 Prometheus 集成，请在 Prometheus 配置文件中添加以下 scrape 配置：

```yaml
scrape_configs:
  - job_name: 'i18n-flow'
    scrape_interval: 15s
    metrics_path: '/metrics'
    static_configs:
      - targets: ['localhost:8080']
```

## 最佳实践

1. **生产环境配置**：在生产环境中，建议设置 IP 白名单以限制对监控指标的访问
2. **告警阈值调整**：根据系统实际负载情况调整告警阈值
3. **监控指标持久化**：配置 Prometheus 数据持久化，以便长期分析系统性能趋势
4. **定期审查仪表盘**：定期检查 Grafana 仪表盘，识别潜在的性能问题
5. **告警通知渠道**：配置多个告警通知渠道，确保关键告警不会被遗漏

## 故障排除

### 无法访问监控指标

- 检查 `METRICS_ENABLED` 是否设置为 `true`
- 检查 IP 白名单配置
- 确认服务器防火墙设置允许访问监控端点

### 指标收集不工作

- 检查日志中是否有与监控系统相关的错误
- 验证 Prometheus 配置是否正确
- 确认服务有足够的资源运行监控系统

### 告警未触发

- 检查 `METRICS_ALERTING_ENABLED` 是否设置为 `true`
- 验证告警阈值配置
- 检查告警通知渠道配置
