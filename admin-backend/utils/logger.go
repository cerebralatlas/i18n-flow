package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 日志配置结构
type LogConfig struct {
	Level      string `json:"level"`       // 日志级别: debug, info, warn, error
	Format     string `json:"format"`      // 日志格式: json, console
	Output     string `json:"output"`      // 输出方式: stdout, file, both
	Filename   string `json:"filename"`    // 日志文件路径
	MaxSize    int    `json:"max_size"`    // 日志文件最大大小(MB)
	MaxAge     int    `json:"max_age"`     // 日志文件保留天数
	MaxBackups int    `json:"max_backups"` // 最大备份文件数
	Compress   bool   `json:"compress"`    // 是否压缩备份文件
}

// Logger 全局日志实例
var Logger *zap.Logger

// SugaredLogger 全局语法糖日志实例
var SugaredLogger *zap.SugaredLogger

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	// 设置日志级别
	level := parseLogLevel(config.Level)

	// 创建编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 根据格式选择编码器
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建写入器
	var cores []zapcore.Core

	// 控制台输出
	if config.Output == "stdout" || config.Output == "both" {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if config.Output == "file" || config.Output == "both" {
		// 确保日志目录存在
		if err := os.MkdirAll(filepath.Dir(config.Filename), 0755); err != nil {
			return err
		}

		fileWriter := &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
		}

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // 文件日志总是使用JSON格式
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 创建日志核心
	core := zapcore.NewTee(cores...)

	// 创建日志实例
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	SugaredLogger = Logger.Sugar()

	return nil
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

// GetDefaultLogConfig 获取默认日志配置
func GetDefaultLogConfig(env string) LogConfig {
	config := LogConfig{
		Level:      "info",
		Format:     "console",
		Output:     "both",
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 5,
		Compress:   true,
	}

	// 根据环境调整配置
	switch env {
	case "development", "dev":
		config.Level = "debug"
		config.Format = "console"
		config.Output = "both"
	case "production", "prod":
		config.Level = "info"
		config.Format = "json"
		config.Output = "file"
	case "test":
		config.Level = "error"
		config.Output = "stdout"
	}

	return config
}

// 便捷日志函数
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Debugf 格式化debug日志
func Debugf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Debugf(template, args...)
	}
}

func Infof(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Infof(template, args...)
	}
}

func Warnf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Warnf(template, args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Errorf(template, args...)
	}
}

func Fatalf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Fatalf(template, args...)
	}
}

// WithFields 创建带字段的日志条目
func WithFields(fields ...zap.Field) *zap.Logger {
	if Logger != nil {
		return Logger.With(fields...)
	}
	return zap.NewNop()
}

// Sync 同步日志缓冲区
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}
