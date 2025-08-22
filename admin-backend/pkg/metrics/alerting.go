package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertManager 告警管理器
type AlertManager struct {
	// 告警配置
	config AlertConfig
	
	// 告警阈值
	thresholds map[string]AlertThreshold
	
	// 告警状态
	alertStates map[string]AlertState
	
	// 告警通知通道
	notifiers []AlertNotifier
	
	// 停止通道
	stopCh chan struct{}
}

// AlertConfig 告警配置
type AlertConfig struct {
	// 检查间隔
	CheckInterval time.Duration
	
	// 恢复通知
	SendRecoveryNotifications bool
	
	// 抑制重复告警的时间窗口
	SuppressWindow time.Duration
}

// AlertThreshold 告警阈值
type AlertThreshold struct {
	// 指标名称
	MetricName string
	
	// 指标标签
	MetricLabels map[string]string
	
	// 阈值
	Threshold float64
	
	// 比较操作符: >, <, >=, <=, ==
	Operator string
	
	// 告警级别: info, warning, error, critical
	Severity string
	
	// 告警描述
	Description string
}

// AlertState 告警状态
type AlertState struct {
	// 是否已触发
	Firing bool
	
	// 触发时间
	FiredAt time.Time
	
	// 最后通知时间
	LastNotifiedAt time.Time
	
	// 恢复时间
	ResolvedAt time.Time
	
	// 当前值
	CurrentValue float64
}

// Alert 告警信息
type Alert struct {
	// 告警名称
	Name string `json:"name"`
	
	// 告警级别
	Severity string `json:"severity"`
	
	// 告警描述
	Description string `json:"description"`
	
	// 告警状态: firing, resolved
	Status string `json:"status"`
	
	// 告警时间
	Timestamp time.Time `json:"timestamp"`
	
	// 指标名称
	MetricName string `json:"metric_name"`
	
	// 指标标签
	MetricLabels map[string]string `json:"metric_labels"`
	
	// 当前值
	CurrentValue float64 `json:"current_value"`
	
	// 阈值
	Threshold float64 `json:"threshold"`
	
	// 比较操作符
	Operator string `json:"operator"`
}

// AlertNotifier 告警通知接口
type AlertNotifier interface {
	// 发送告警通知
	Notify(alert Alert) error
}

// WebhookNotifier Webhook告警通知器
type WebhookNotifier struct {
	URL string
}

// EmailNotifier 邮件告警通知器
type EmailNotifier struct {
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	ToEmails     []string
}

// NewAlertManager 创建告警管理器
func NewAlertManager(config AlertConfig) *AlertManager {
	return &AlertManager{
		config:      config,
		thresholds:  make(map[string]AlertThreshold),
		alertStates: make(map[string]AlertState),
		notifiers:   []AlertNotifier{},
		stopCh:      make(chan struct{}),
	}
}

// AddThreshold 添加告警阈值
func (am *AlertManager) AddThreshold(name string, threshold AlertThreshold) {
	am.thresholds[name] = threshold
}

// AddNotifier 添加告警通知器
func (am *AlertManager) AddNotifier(notifier AlertNotifier) {
	am.notifiers = append(am.notifiers, notifier)
}

// Start 开始告警检查
func (am *AlertManager) Start() {
	go func() {
		ticker := time.NewTicker(am.config.CheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				am.checkAlerts()
			case <-am.stopCh:
				return
			}
		}
	}()
}

// Stop 停止告警检查
func (am *AlertManager) Stop() {
	close(am.stopCh)
}

// checkAlerts 检查所有告警
func (am *AlertManager) checkAlerts() {
	for name, threshold := range am.thresholds {
		am.checkAlert(name, threshold)
	}
}

// checkAlert 检查单个告警
func (am *AlertManager) checkAlert(name string, threshold AlertThreshold) {
	// 获取当前指标值
	value := getMetricValue(threshold.MetricName, threshold.MetricLabels)
	
	// 获取当前告警状态
	state, exists := am.alertStates[name]
	if !exists {
		state = AlertState{}
	}
	
	// 更新当前值
	state.CurrentValue = value
	
	// 检查是否触发告警
	firing := evaluateThreshold(value, threshold.Threshold, threshold.Operator)
	
	// 处理告警状态变化
	if firing && !state.Firing {
		// 新触发的告警
		state.Firing = true
		state.FiredAt = time.Now()
		state.LastNotifiedAt = time.Now()
		
		// 发送告警通知
		am.sendAlert(name, threshold, state, "firing")
	} else if firing && state.Firing {
		// 持续触发的告警
		now := time.Now()
		if now.Sub(state.LastNotifiedAt) > am.config.SuppressWindow {
			// 超过抑制窗口，再次发送通知
			state.LastNotifiedAt = now
			am.sendAlert(name, threshold, state, "firing")
		}
	} else if !firing && state.Firing {
		// 告警已恢复
		state.Firing = false
		state.ResolvedAt = time.Now()
		
		// 发送恢复通知
		if am.config.SendRecoveryNotifications {
			am.sendAlert(name, threshold, state, "resolved")
		}
	}
	
	// 更新告警状态
	am.alertStates[name] = state
}

// evaluateThreshold 评估阈值
func evaluateThreshold(value, threshold float64, operator string) bool {
	switch operator {
	case ">":
		return value > threshold
	case "<":
		return value < threshold
	case ">=":
		return value >= threshold
	case "<=":
		return value <= threshold
	case "==":
		return value == threshold
	default:
		return false
	}
}

// getMetricValue 获取指标值
func getMetricValue(metricName string, labels map[string]string) float64 {
	// 这里应该实现从Prometheus查询指标值的逻辑
	// 由于这需要与Prometheus API交互，这里只是一个占位符
	return 0
}

// sendAlert 发送告警通知
func (am *AlertManager) sendAlert(name string, threshold AlertThreshold, state AlertState, status string) {
	alert := Alert{
		Name:         name,
		Severity:     threshold.Severity,
		Description:  threshold.Description,
		Status:       status,
		Timestamp:    time.Now(),
		MetricName:   threshold.MetricName,
		MetricLabels: threshold.MetricLabels,
		CurrentValue: state.CurrentValue,
		Threshold:    threshold.Threshold,
		Operator:     threshold.Operator,
	}
	
	// 发送到所有通知器
	for _, notifier := range am.notifiers {
		notifier.Notify(alert)
	}
}

// Notify 实现WebhookNotifier的Notify方法
func (n *WebhookNotifier) Notify(alert Alert) error {
	// 将告警信息转换为JSON
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	
	// 发送HTTP POST请求
	resp, err := http.Post(n.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d", resp.StatusCode)
	}
	
	return nil
}

// Notify 实现EmailNotifier的Notify方法
func (n *EmailNotifier) Notify(alert Alert) error {
	// 这里应该实现发送邮件的逻辑
	// 由于这需要SMTP客户端，这里只是一个占位符
	return nil
}
