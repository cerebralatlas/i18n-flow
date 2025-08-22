package utils_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.uber.org/zap/zaptest/observer"

	"i18n-flow/utils"
)

func TestParseLogLevel(t *testing.T) {
	// 由于parseLogLevel是私有函数，我们通过测试GetDefaultLogConfig来间接测试
	tests := []struct {
		env           string
		expectedLevel string
	}{
		{"development", "debug"},
		{"dev", "debug"},
		{"production", "info"},
		{"prod", "info"},
		{"test", "error"},
		{"unknown", "info"}, // 默认应该是info
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := utils.GetDefaultLogConfig(tt.env)
			assert.Equal(t, tt.expectedLevel, config.Level)
		})
	}
}

func TestGetDefaultLogConfig(t *testing.T) {
	tests := []struct {
		env      string
		expected utils.LogConfig
	}{
		{
			"development",
			utils.LogConfig{
				Level:      "debug",
				Format:     "console",
				Output:     "both",
				Filename:   "logs/app.log",
				MaxSize:    100,
				MaxAge:     7,
				MaxBackups: 5,
				Compress:   true,
			},
		},
		{
			"production",
			utils.LogConfig{
				Level:      "info",
				Format:     "json",
				Output:     "file",
				Filename:   "logs/app.log",
				MaxSize:    100,
				MaxAge:     7,
				MaxBackups: 5,
				Compress:   true,
			},
		},
		{
			"test",
			utils.LogConfig{
				Level:      "error",
				Format:     "console",
				Output:     "stdout",
				Filename:   "logs/app.log",
				MaxSize:    100,
				MaxAge:     7,
				MaxBackups: 5,
				Compress:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := utils.GetDefaultLogConfig(tt.env)
			assert.Equal(t, tt.expected.Level, config.Level)
			assert.Equal(t, tt.expected.Format, config.Format)
			assert.Equal(t, tt.expected.Output, config.Output)
			assert.Equal(t, tt.expected.Filename, config.Filename)
		})
	}
}

func TestInitLogger(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := ioutil.TempDir("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name   string
		config utils.LogConfig
	}{
		{
			"Console Output",
			utils.LogConfig{
				Level:      "debug",
				Format:     "console",
				Output:     "stdout",
				Filename:   filepath.Join(tempDir, "app.log"),
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 1,
				Compress:   false,
			},
		},
		{
			"File Output",
			utils.LogConfig{
				Level:      "info",
				Format:     "json",
				Output:     "file",
				Filename:   filepath.Join(tempDir, "app.log"),
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 1,
				Compress:   false,
			},
		},
		{
			"Both Outputs",
			utils.LogConfig{
				Level:      "warn",
				Format:     "json",
				Output:     "both",
				Filename:   filepath.Join(tempDir, "app.log"),
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 1,
				Compress:   false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化日志
			err := utils.InitLogger(tc.config)
			require.NoError(t, err)

			// 验证Logger和SugaredLogger已经初始化
			assert.NotNil(t, utils.Logger)
			assert.NotNil(t, utils.SugaredLogger)

			// 测试日志记录函数
			utils.Debug("debug message")
			utils.Info("info message")
			utils.Warn("warning message")
			utils.Error("error message")
			// 不测试Fatal，因为它会导致程序退出

			// 测试格式化日志函数
			utils.Debugf("debug %s", "formatted")
			utils.Infof("info %s", "formatted")
			utils.Warnf("warn %s", "formatted")
			utils.Errorf("error %s", "formatted")

			// 测试WithFields
			logger := utils.WithFields(zap.String("key", "value"))
			assert.NotNil(t, logger)

			// 同步日志
			utils.Sync()

			// 如果是文件输出，验证文件是否存在
			if tc.config.Output == "file" || tc.config.Output == "both" {
				_, err := os.Stat(tc.config.Filename)
				assert.NoError(t, err)
			}
		})
	}
}

// 测试日志记录函数的行为
func TestLogFunctions(t *testing.T) {
	// 使用zaptest创建一个测试观察者
	core, logs := observer.New(zap.InfoLevel)

	// 临时替换全局Logger
	originalLogger := utils.Logger
	utils.Logger = zap.New(core)
	utils.SugaredLogger = utils.Logger.Sugar()
	defer func() {
		utils.Logger = originalLogger
		utils.SugaredLogger = originalLogger.Sugar()
	}()

	// 测试各种日志级别
	utils.Info("test info")
	utils.Warn("test warn")
	utils.Error("test error")

	// 验证日志条目
	assert.Equal(t, 3, logs.Len())
	assert.Equal(t, "test info", logs.All()[0].Message)
	assert.Equal(t, "test warn", logs.All()[1].Message)
	assert.Equal(t, "test error", logs.All()[2].Message)
}

// 测试在Logger为nil时的行为
func TestNilLogger(t *testing.T) {
	// 临时保存并清除全局Logger
	originalLogger := utils.Logger
	originalSugar := utils.SugaredLogger
	utils.Logger = nil
	utils.SugaredLogger = nil
	defer func() {
		utils.Logger = originalLogger
		utils.SugaredLogger = originalSugar
	}()

	// 这些调用不应该导致panic
	utils.Debug("debug")
	utils.Info("info")
	utils.Warn("warn")
	utils.Error("error")
	utils.Debugf("debug %s", "format")
	utils.Infof("info %s", "format")
	utils.Warnf("warn %s", "format")
	utils.Errorf("error %s", "format")
	logger := utils.WithFields(zap.String("key", "value"))
	assert.NotNil(t, logger) // 应该返回一个nop logger
	utils.Sync()             // 不应该panic
}
