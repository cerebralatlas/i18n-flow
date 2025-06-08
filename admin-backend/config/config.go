package config

import (
	"log"
	"os"
	"strconv"

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
	DB  DBConfig
	JWT JWTConfig
	CLI CLIConfig
	Log LogConfig
}

// GetConfig 获取配置
func GetConfig() *Config {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: .env文件未找到，将使用默认配置或环境变量")
	}

	// 读取端口并转换为整数
	port, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		log.Println("警告: 无法解析数据库端口，使用默认值3306")
		port = 3306
	}

	// 解析JWT过期时间
	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		jwtExpHours = 24
	}

	jwtRefreshExpHours, err := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION_HOURS", "168"))
	if err != nil {
		jwtRefreshExpHours = 168 // 7天
	}

	// 解析日志配置
	logMaxSize, err := strconv.Atoi(getEnv("LOG_MAX_SIZE", "100"))
	if err != nil {
		logMaxSize = 100
	}

	logMaxAge, err := strconv.Atoi(getEnv("LOG_MAX_AGE", "7"))
	if err != nil {
		logMaxAge = 7
	}

	logMaxBackups, err := strconv.Atoi(getEnv("LOG_MAX_BACKUPS", "5"))
	if err != nil {
		logMaxBackups = 5
	}

	logCompress := getEnv("LOG_COMPRESS", "true") == "true"

	return &Config{
		DB: DBConfig{
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     port,
			DBName:   getEnv("DB_NAME", "i18n_flow"),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "your-256-bit-secret"),
			ExpirationHours:        jwtExpHours,
			RefreshSecret:          getEnv("JWT_REFRESH_SECRET", "your-refresh-secret"),
			RefreshExpirationHours: jwtRefreshExpHours,
		},
		CLI: CLIConfig{
			APIKey: getEnv("CLI_API_KEY", "testapikey"),
		},
		Log: LogConfig{
			Level:         getEnv("LOG_LEVEL", "info"),
			Format:        getEnv("LOG_FORMAT", "console"),
			Output:        getEnv("LOG_OUTPUT", "both"),
			LogDir:        getEnv("LOG_DIR", "logs"),
			DateFormat:    getEnv("LOG_DATE_FORMAT", "2006-01-02"),
			MaxSize:       logMaxSize,
			MaxAge:        logMaxAge,
			MaxBackups:    logMaxBackups,
			Compress:      logCompress,
			EnableConsole: getEnv("LOG_ENABLE_CONSOLE", "true") == "true",
		},
	}
}

// getEnv 从环境变量获取值，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
