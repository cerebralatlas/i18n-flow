package config

// MetricsConfig 监控配置
type MetricsConfig struct {
	// 是否启用监控
	Enabled bool `json:"enabled" yaml:"enabled"`
	
	// 允许访问监控指标的IP白名单
	AllowedIPs []string `json:"allowed_ips" yaml:"allowed_ips"`
	
	// 是否启用系统指标收集
	SystemMetricsEnabled bool `json:"system_metrics_enabled" yaml:"system_metrics_enabled"`
	
	// 是否启用数据库指标收集
	DBMetricsEnabled bool `json:"db_metrics_enabled" yaml:"db_metrics_enabled"`
	
	// 是否启用Redis指标收集
	RedisMetricsEnabled bool `json:"redis_metrics_enabled" yaml:"redis_metrics_enabled"`
	
	// 是否启用业务指标收集
	BusinessMetricsEnabled bool `json:"business_metrics_enabled" yaml:"business_metrics_enabled"`
	
	// 是否启用告警系统
	AlertingEnabled bool `json:"alerting_enabled" yaml:"alerting_enabled"`
	
	// 告警通知Webhook URL
	AlertWebhookURL string `json:"alert_webhook_url" yaml:"alert_webhook_url"`
	
	// 告警通知邮箱配置
	AlertEmail AlertEmailConfig `json:"alert_email" yaml:"alert_email"`
	
	// 指标收集间隔（秒）
	CollectionIntervalSeconds int `json:"collection_interval_seconds" yaml:"collection_interval_seconds"`
}

// AlertEmailConfig 告警邮箱配置
type AlertEmailConfig struct {
	// 是否启用邮箱告警
	Enabled bool `json:"enabled" yaml:"enabled"`
	
	// SMTP服务器地址
	SMTPServer string `json:"smtp_server" yaml:"smtp_server"`
	
	// SMTP服务器端口
	SMTPPort int `json:"smtp_port" yaml:"smtp_port"`
	
	// SMTP用户名
	SMTPUser string `json:"smtp_user" yaml:"smtp_user"`
	
	// SMTP密码
	SMTPPassword string `json:"smtp_password" yaml:"smtp_password"`
	
	// 发件人邮箱
	FromEmail string `json:"from_email" yaml:"from_email"`
	
	// 收件人邮箱列表
	ToEmails []string `json:"to_emails" yaml:"to_emails"`
}

// GetDefaultMetricsConfig 获取默认监控配置
func GetDefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:                true,
		AllowedIPs:             []string{"127.0.0.1", "::1"},
		SystemMetricsEnabled:   true,
		DBMetricsEnabled:       true,
		RedisMetricsEnabled:    true,
		BusinessMetricsEnabled: true,
		AlertingEnabled:        false,
		AlertWebhookURL:        "",
		AlertEmail: AlertEmailConfig{
			Enabled:      false,
			SMTPServer:   "smtp.example.com",
			SMTPPort:     587,
			SMTPUser:     "user@example.com",
			SMTPPassword: "",
			FromEmail:    "alerts@example.com",
			ToEmails:     []string{"admin@example.com"},
		},
		CollectionIntervalSeconds: 30,
	}
}
