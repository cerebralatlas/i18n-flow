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

// Config 应用配置
type Config struct {
	DB  DBConfig
	JWT JWTConfig
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
