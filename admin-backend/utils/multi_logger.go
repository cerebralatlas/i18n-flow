package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerType 日志类型
type LoggerType string

const (
	LogTypeApp    LoggerType = "app"    // 应用日志
	LogTypeAccess LoggerType = "access" // 访问日志
	LogTypeError  LoggerType = "error"  // 错误日志
	LogTypeAuth   LoggerType = "auth"   // 认证日志
	LogTypeDB     LoggerType = "db"     // 数据库日志
)

// MultiLogConfig 多日志配置
type MultiLogConfig struct {
	Level         string            `json:"level"`          // 全局日志级别
	Format        string            `json:"format"`         // 日志格式: json, console
	Output        string            `json:"output"`         // 输出方式: stdout, file, both
	LogDir        string            `json:"log_dir"`        // 日志目录
	DateFormat    string            `json:"date_format"`    // 日期格式，用于文件名
	MaxSize       int               `json:"max_size"`       // 每个日志文件最大大小(MB)
	MaxAge        int               `json:"max_age"`        // 日志文件保留天数
	MaxBackups    int               `json:"max_backups"`    // 最大备份文件数
	Compress      bool              `json:"compress"`       // 是否压缩备份文件
	EnableConsole bool              `json:"enable_console"` // 是否启用控制台输出
	LogTypes      map[string]string `json:"log_types"`      // 特定类型的日志级别
}

// LoggerManager 日志管理器
type LoggerManager struct {
	config  MultiLogConfig
	loggers map[LoggerType]*zap.Logger
}

var logManager *LoggerManager

// Logger 全局日志实例
var Logger *zap.Logger

// SugaredLogger 全局语法糖日志实例
var SugaredLogger *zap.SugaredLogger

// InitMultiLogger 初始化多日志系统
func InitMultiLogger(config MultiLogConfig) error {
	logManager = &LoggerManager{
		config:  config,
		loggers: make(map[LoggerType]*zap.Logger),
	}

	// 确保日志目录存在
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建不同类型的日志器
	if err := logManager.createLoggers(); err != nil {
		return err
	}

	// 设置全局日志器为应用日志器
	Logger = logManager.GetLogger(LogTypeApp)
	SugaredLogger = Logger.Sugar()

	return nil
}

// createLoggers 创建各种类型的日志器
func (lm *LoggerManager) createLoggers() error {
	logTypes := []LoggerType{LogTypeApp, LogTypeAccess, LogTypeError, LogTypeAuth, LogTypeDB}

	for _, logType := range logTypes {
		logger, err := lm.createLogger(logType)
		if err != nil {
			return fmt.Errorf("创建%s日志器失败: %v", logType, err)
		}
		lm.loggers[logType] = logger
	}

	return nil
}

// createLogger 创建单个日志器
func (lm *LoggerManager) createLogger(logType LoggerType) (*zap.Logger, error) {
	// 根据日志类型确定级别
	level := lm.getLogLevel(logType)

	// 创建编码器配置
	encoderConfig := lm.getEncoderConfig()

	var cores []zapcore.Core

	// 控制台输出
	if lm.config.EnableConsole && (lm.config.Output == "stdout" || lm.config.Output == "both") {
		consoleEncoder := lm.getConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if lm.config.Output == "file" || lm.config.Output == "both" {
		// 为每种日志类型创建单独的文件写入器
		filename := lm.getLogFilename(logType)
		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    lm.config.MaxSize,
			MaxAge:     lm.config.MaxAge,
			MaxBackups: lm.config.MaxBackups,
			Compress:   lm.config.Compress,
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 错误日志额外写入error.log
	if logType != LogTypeError && level <= zapcore.ErrorLevel {
		errorFilename := lm.getLogFilename(LogTypeError)
		errorWriter := &lumberjack.Logger{
			Filename:   errorFilename,
			MaxSize:    lm.config.MaxSize,
			MaxAge:     lm.config.MaxAge,
			MaxBackups: lm.config.MaxBackups,
			Compress:   lm.config.Compress,
		}

		errorCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(errorWriter),
			zapcore.ErrorLevel,
		)
		cores = append(cores, errorCore)
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// getLogLevel 获取指定类型的日志级别
func (lm *LoggerManager) getLogLevel(logType LoggerType) zapcore.Level {
	// 检查是否有特定配置
	if levelStr, exists := lm.config.LogTypes[string(logType)]; exists {
		return parseLogLevel(levelStr)
	}

	// 使用全局配置
	return parseLogLevel(lm.config.Level)
}

// getEncoderConfig 获取编码器配置
func (lm *LoggerManager) getEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.LevelKey = "level"
	config.NameKey = "logger"
	config.CallerKey = "caller"
	config.MessageKey = "message"
	config.StacktraceKey = "stacktrace"
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	return config
}

// getConsoleEncoder 获取控制台编码器
func (lm *LoggerManager) getConsoleEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	if lm.config.Format == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	// 控制台使用彩色输出
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}

// getLogFilename 获取日志文件名
func (lm *LoggerManager) getLogFilename(logType LoggerType) string {
	dateStr := time.Now().Format(lm.config.DateFormat)
	filename := fmt.Sprintf("%s-%s.log", logType, dateStr)
	return filepath.Join(lm.config.LogDir, filename)
}

// GetLogger 获取指定类型的日志器
func (lm *LoggerManager) GetLogger(logType LoggerType) *zap.Logger {
	if logger, exists := lm.loggers[logType]; exists {
		return logger
	}
	return zap.NewNop()
}

// 全局访问函数
func GetLogger(logType LoggerType) *zap.Logger {
	if logManager != nil {
		return logManager.GetLogger(logType)
	}
	return zap.NewNop()
}

// 不同类型的日志记录函数
// 全局日志函数（兼容旧代码）
func Debug(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Fatal(msg, fields...)
}

// 应用日志函数
func AppInfo(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Info(msg, fields...)
}

func AppError(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Error(msg, fields...)
}

func AppWarn(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Warn(msg, fields...)
}

func AppDebug(msg string, fields ...zap.Field) {
	GetLogger(LogTypeApp).Debug(msg, fields...)
}

func AccessLog(msg string, fields ...zap.Field) {
	GetLogger(LogTypeAccess).Info(msg, fields...)
}

func ErrorLog(msg string, fields ...zap.Field) {
	GetLogger(LogTypeError).Error(msg, fields...)
}

func AuthLog(msg string, fields ...zap.Field) {
	GetLogger(LogTypeAuth).Info(msg, fields...)
}

func DBLog(msg string, fields ...zap.Field) {
	GetLogger(LogTypeDB).Info(msg, fields...)
}

// 格式化日志函数（兼容旧代码）
func Debugf(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Fatalf(template, args...)
}

// 应用格式化日志函数
func AppInfof(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Infof(template, args...)
}

func AppErrorf(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Errorf(template, args...)
}

func AppWarnf(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Warnf(template, args...)
}

func AppDebugf(template string, args ...interface{}) {
	GetLogger(LogTypeApp).Sugar().Debugf(template, args...)
}

func AccessLogf(template string, args ...interface{}) {
	GetLogger(LogTypeAccess).Sugar().Infof(template, args...)
}

// SyncAll 同步所有日志缓冲区
func SyncAll() {
	if logManager != nil {
		for _, logger := range logManager.loggers {
			logger.Sync()
		}
	}
}

// Sync 同步日志缓冲区（兼容旧代码）
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

// WithFields 创建带字段的日志条目（兼容旧代码）
func WithFields(fields ...zap.Field) *zap.Logger {
	if Logger != nil {
		return Logger.With(fields...)
	}
	return zap.NewNop()
}

// GetDefaultMultiLogConfig 获取默认多日志配置
func GetDefaultMultiLogConfig(env string) MultiLogConfig {
	config := MultiLogConfig{
		Level:         "info",
		Format:        "console",
		Output:        "both",
		LogDir:        "logs",
		DateFormat:    "2006-01-02",
		MaxSize:       100,
		MaxAge:        30,
		MaxBackups:    10,
		Compress:      true,
		EnableConsole: true,
		LogTypes: map[string]string{
			"access": "info",
			"error":  "error",
			"auth":   "info",
			"db":     "warn",
		},
	}

	// 根据环境调整配置
	switch env {
	case "development", "dev":
		config.Level = "debug"
		config.Format = "console"
		config.Output = "both"
		config.EnableConsole = true
		config.LogTypes["db"] = "debug"
	case "production", "prod":
		config.Level = "info"
		config.Format = "json"
		config.Output = "file"
		config.EnableConsole = false
		config.MaxAge = 60
	case "test":
		config.Level = "error"
		config.Output = "stdout"
		config.EnableConsole = true
	}

	return config
}
