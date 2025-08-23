package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init 初始化日志
// mode: "dev" | "prod"
func Init(mode string) {
	var cfg zap.Config
	if mode == "dev" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	l, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	log = l
}

// Sync 确保日志刷到文件或控制台
func Sync() {
	_ = log.Sync()
}

// 常用方法
func Debug(msg string, fields ...zap.Field) { log.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { log.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { log.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { log.Error(msg, fields...) }

// WithFields 生成带字段的 logger（链式调用）
func WithFields(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// ExampleField 预定义字段（比如 requestID、userID）
func String(key, val string) zap.Field  { return zap.String(key, val) }
func Int(key string, val int) zap.Field { return zap.Int(key, val) }
