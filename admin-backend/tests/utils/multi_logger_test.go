package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"i18n-flow/utils"
)

func TestGetDefaultMultiLogConfig(t *testing.T) {
	tests := []struct {
		env             string
		expectedLevel   string
		expectedFormat  string
		expectedOutput  string
		expectedConsole bool
		expectedDBLevel string
	}{
		{
			"development",
			"debug",
			"console",
			"both",
			true,
			"debug",
		},
		{
			"dev",
			"debug",
			"console",
			"both",
			true,
			"debug",
		},
		{
			"production",
			"info",
			"json",
			"file",
			false,
			"warn", // 默认值
		},
		{
			"prod",
			"info",
			"json",
			"file",
			false,
			"warn", // 默认值
		},
		{
			"test",
			"error",
			"console",
			"stdout",
			true,
			"warn", // 默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := utils.GetDefaultMultiLogConfig(tt.env)
			assert.Equal(t, tt.expectedLevel, config.Level)
			assert.Equal(t, tt.expectedFormat, config.Format)
			assert.Equal(t, tt.expectedOutput, config.Output)
			assert.Equal(t, tt.expectedConsole, config.EnableConsole)
			assert.Equal(t, tt.expectedDBLevel, config.LogTypes["db"])
		})
	}
}

func TestInitMultiLogger(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "multi_logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试配置
	config := utils.MultiLogConfig{
		Level:         "info",
		Format:        "json",
		Output:        "both",
		LogDir:        tempDir,
		DateFormat:    "2006-01-02",
		MaxSize:       1,
		MaxAge:        1,
		MaxBackups:    1,
		Compress:      false,
		EnableConsole: true,
		LogTypes: map[string]string{
			"access": "info",
			"error":  "error",
			"auth":   "info",
			"db":     "warn",
		},
	}

	// 初始化多日志系统
	err = utils.InitMultiLogger(config)
	require.NoError(t, err)

	// 验证Logger和SugaredLogger已经初始化
	assert.NotNil(t, utils.Logger)
	assert.NotNil(t, utils.SugaredLogger)

	// 测试不同类型的日志记录函数
	utils.AppInfo("app info message")
	utils.AppError("app error message")
	utils.AppWarn("app warning message")
	utils.AppDebug("app debug message")
	utils.AccessLog("access log message")
	utils.ErrorLog("error log message")
	utils.AuthLog("auth log message")
	utils.DBLog("db log message")

	// 测试格式化日志函数
	utils.AppInfof("app info %s", "formatted")
	utils.AppErrorf("app error %s", "formatted")
	utils.AccessLogf("access log %s", "formatted")

	// 同步所有日志
	utils.SyncAll()

	// 验证日志文件是否创建
	// 由于日志文件名包含日期，我们只检查目录中是否有文件
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	assert.NotEmpty(t, files)
}

func TestGetLogger(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "get_logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试配置
	config := utils.MultiLogConfig{
		Level:         "info",
		Format:        "console",
		Output:        "stdout", // 仅控制台输出，避免创建文件
		LogDir:        tempDir,
		DateFormat:    "2006-01-02",
		MaxSize:       1,
		MaxAge:        1,
		MaxBackups:    1,
		Compress:      false,
		EnableConsole: true,
		LogTypes: map[string]string{
			"access": "info",
			"error":  "error",
			"auth":   "info",
			"db":     "warn",
		},
	}

	// 初始化多日志系统
	err = utils.InitMultiLogger(config)
	require.NoError(t, err)

	// 测试获取不同类型的日志器
	logTypes := []utils.LoggerType{
		utils.LogTypeApp,
		utils.LogTypeAccess,
		utils.LogTypeError,
		utils.LogTypeAuth,
		utils.LogTypeDB,
	}

	for _, logType := range logTypes {
		logger := utils.GetLogger(logType)
		assert.NotNil(t, logger)
	}

	// 测试获取不存在的日志类型
	unknownLogger := utils.GetLogger("unknown")
	assert.NotNil(t, unknownLogger) // 应该返回一个nop logger
}

func TestLoggerFunctions(t *testing.T) {
	// 使用zaptest创建一个测试观察者
	core, logs := observer.New(zap.InfoLevel)

	// 保存原始日志器并在测试后恢复
	originalLogger := utils.Logger
	originalSugar := utils.SugaredLogger

	// 设置测试日志器
	utils.Logger = zap.New(core)
	utils.SugaredLogger = utils.Logger.Sugar()

	defer func() {
		utils.Logger = originalLogger
		utils.SugaredLogger = originalSugar
	}()

	// 测试日志函数
	utils.Info("test app info")

	// 验证日志条目
	assert.Equal(t, 1, logs.Len(), "Expected 1 log entry but got %d", logs.Len())
	if logs.Len() > 0 {
		assert.Equal(t, "test app info", logs.All()[0].Message)
	}
}

// 测试在logManager为nil时的行为
func TestNilLogManager(t *testing.T) {
	// 这些调用不应该导致panic
	logger := utils.GetLogger(utils.LogTypeApp)
	assert.NotNil(t, logger) // 应该返回一个nop logger

	// 测试日志函数，不应该panic
	utils.AppInfo("test info")
	utils.AppError("test error")
	utils.AppWarn("test warn")
	utils.AppDebug("test debug")
	utils.AccessLog("test access")
	utils.ErrorLog("test error log")
	utils.AuthLog("test auth")
	utils.DBLog("test db")

	utils.AppInfof("test %s", "infof")
	utils.AppErrorf("test %s", "errorf")
	utils.AccessLogf("test %s", "accessf")

	// 同步不应该panic
	utils.SyncAll()
}
