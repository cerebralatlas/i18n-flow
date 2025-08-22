package metrics

import (
	"runtime"
	"time"
	
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMetricsCollector 系统指标收集器
type SystemMetricsCollector struct {
	interval time.Duration
	stopCh   chan struct{}
}

// NewSystemMetricsCollector 创建系统指标收集器
func NewSystemMetricsCollector(interval time.Duration) *SystemMetricsCollector {
	return &SystemMetricsCollector{
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start 开始收集系统指标
func (c *SystemMetricsCollector) Start() {
	// 立即收集一次指标
	c.collectMetrics()
	
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

// Stop 停止收集系统指标
func (c *SystemMetricsCollector) Stop() {
	close(c.stopCh)
}

// collectMetrics 收集系统指标
func (c *SystemMetricsCollector) collectMetrics() {
	// 收集内存指标
	c.collectMemoryMetrics()
	
	// 收集CPU指标
	c.collectCPUMetrics()
}

// collectMemoryMetrics 收集内存指标
func (c *SystemMetricsCollector) collectMemoryMetrics() {
	// 收集Go运行时内存统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// 设置Go堆内存使用量
	SetMemoryUsage(float64(m.HeapAlloc))

	// 使用gopsutil收集系统内存信息
	if v, err := mem.VirtualMemory(); err == nil {
		// 这里可以添加更多系统内存指标
		_ = v // 避免未使用变量警告
	}
}

// collectCPUMetrics 收集CPU指标
func (c *SystemMetricsCollector) collectCPUMetrics() {
	// 使用gopsutil收集CPU信息
	if percent, err := cpu.Percent(time.Second, false); err == nil && len(percent) > 0 {
		SetCPUUsage(percent[0])
	}
}