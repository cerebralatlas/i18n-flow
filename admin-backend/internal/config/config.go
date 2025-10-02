package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// DBConfig 数据库配置
type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret                 string
	ExpirationHours        int
	RefreshSecret          string
	RefreshExpirationHours int
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Prefix   string
}

// CLIConfig CLI配置
type CLIConfig struct {
	APIKey string
}

// LogConfig 日志配置
type LogConfig struct {
	Level         string `json:"level"`          // 全局日志级别
	Format        string `json:"format"`         // 日志格式
	Output        string `json:"output"`         // 输出方式
	LogDir        string `json:"log_dir"`        // 日志目录
	DateFormat    string `json:"date_format"`    // 日期格式
	MaxSize       int    `json:"max_size"`       // 最大文件大小
	MaxAge        int    `json:"max_age"`        // 保留天数
	MaxBackups    int    `json:"max_backups"`    // 最大备份数
	Compress      bool   `json:"compress"`       // 是否压缩
	EnableConsole bool   `json:"enable_console"` // 是否启用控制台
}

// Config 应用配置
type Config struct {
	DB    DBConfig
	JWT   JWTConfig
	CLI   CLIConfig
	Log   LogConfig
	Redis RedisConfig
}

// Load 加载配置
func Load() (*Config, error) {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: .env文件未找到，将使用默认配置或环境变量")
	}

	config := &Config{
		DB: DBConfig{
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			DBName:   getEnv("DB_NAME", "i18n_flow"),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "your-256-bit-secret"),
			ExpirationHours:        getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshSecret:          getEnv("JWT_REFRESH_SECRET", "your-refresh-secret"),
			RefreshExpirationHours: getEnvAsInt("JWT_REFRESH_EXPIRATION_HOURS", 168),
		},
		CLI: CLIConfig{
			APIKey: getEnv("CLI_API_KEY", "testapikey"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Prefix:   getEnv("REDIS_PREFIX", "i18n_flow:"),
		},
		Log: LogConfig{
			Level:         getEnv("LOG_LEVEL", "info"),
			Format:        getEnv("LOG_FORMAT", "console"),
			Output:        getEnv("LOG_OUTPUT", "both"),
			LogDir:        getEnv("LOG_DIR", "logs"),
			DateFormat:    getEnv("LOG_DATE_FORMAT", "2006-01-02"),
			MaxSize:       getEnvAsInt("LOG_MAX_SIZE", 100),
			MaxAge:        getEnvAsInt("LOG_MAX_AGE", 7),
			MaxBackups:    getEnvAsInt("LOG_MAX_BACKUPS", 5),
			Compress:      getEnvAsBool("LOG_COMPRESS", true),
			EnableConsole: getEnvAsBool("LOG_ENABLE_CONSOLE", true),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// JWT配置验证 - 强化安全检查
	if c.JWT.Secret == "" {
		return errors.New("JWT secret must be set")
	}
	if c.JWT.Secret == "your-256-bit-secret" || c.JWT.Secret == "your-256-bit-secret-change-this-to-a-secure-random-string" {
		return errors.New("JWT secret must not use default value - please generate a secure random key")
	}
	if len(c.JWT.Secret) < 32 {
		return errors.New("JWT secret must be at least 32 characters long")
	}
	// 检查密钥复杂度，确保包含足够的熵
	if !isStrongKey(c.JWT.Secret) {
		return errors.New("JWT secret must contain a mix of uppercase, lowercase, numbers and special characters")
	}

	if c.JWT.RefreshSecret == "" {
		return errors.New("JWT refresh secret must be set")
	}
	if c.JWT.RefreshSecret == "your-refresh-secret" || c.JWT.RefreshSecret == "your-refresh-secret-change-this-to-a-secure-random-string" {
		return errors.New("JWT refresh secret must not use default value - please generate a secure random key")
	}
	if len(c.JWT.RefreshSecret) < 32 {
		return errors.New("JWT refresh secret must be at least 32 characters long")
	}
	// 刷新密钥必须与主密钥不同
	if c.JWT.RefreshSecret == c.JWT.Secret {
		return errors.New("JWT refresh secret must be different from the main secret")
	}
	if !isStrongKey(c.JWT.RefreshSecret) {
		return errors.New("JWT refresh secret must contain a mix of uppercase, lowercase, numbers and special characters")
	}

	if c.JWT.ExpirationHours <= 0 || c.JWT.ExpirationHours > 24*7 {
		return errors.New("JWT expiration hours must be between 1 and 168 (7 days)")
	}

	if c.JWT.RefreshExpirationHours <= 0 || c.JWT.RefreshExpirationHours > 24*30 {
		return errors.New("JWT refresh expiration hours must be between 1 and 720 (30 days)")
	}

	// 数据库配置验证
	if c.DB.Username == "" {
		return errors.New("database username is required")
	}

	if c.DB.DBName == "" {
		return errors.New("database name is required")
	}

	if c.DB.Host == "" {
		return errors.New("database host is required")
	}

	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	// CLI配置验证
	if c.CLI.APIKey == "" || c.CLI.APIKey == "testapikey" {
		return errors.New("CLI API key must be set and not use default value")
	}
	if len(c.CLI.APIKey) < 16 {
		return errors.New("CLI API key must be at least 16 characters long")
	}

	// Redis配置验证
	if c.Redis.Host == "" {
		return errors.New("Redis host is required")
	}

	if c.Redis.Port <= 0 || c.Redis.Port > 65535 {
		return errors.New("Redis port must be between 1 and 65535")
	}

	if c.Redis.DB < 0 || c.Redis.DB > 15 {
		return errors.New("Redis DB must be between 0 and 15")
	}

	// 日志配置验证
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLogLevels[c.Log.Level] {
		return errors.New("log level must be one of: debug, info, warn, error, fatal")
	}

	if c.Log.MaxSize <= 0 || c.Log.MaxSize > 1000 {
		return errors.New("log max size must be between 1 and 1000 MB")
	}

	if c.Log.MaxAge <= 0 || c.Log.MaxAge > 365 {
		return errors.New("log max age must be between 1 and 365 days")
	}

	if c.Log.MaxBackups < 0 || c.Log.MaxBackups > 100 {
		return errors.New("log max backups must be between 0 and 100")
	}

	return nil
}

// GetConfig 获取配置
func GetConfig() (*Config, error) {
	config, err := Load()
	if err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}
	return config, nil
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(getEnv(key, ""))
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}

// isStrongKey 检查密钥强度
func isStrongKey(key string) bool {
	if len(key) < 32 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range key {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?~", char):
			hasSpecial = true
		}
	}

	// 至少包含3种字符类型
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}
