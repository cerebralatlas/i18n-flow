package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv           string
	AppPort          string
	MongoURI         string
	MongoDB          string
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpire  int
	JWTRefreshExpire int
}

var Cfg *Config

// Init 初始化配置
func Init() {
	viper.SetConfigFile(".env") // 指定 .env 文件
	viper.AutomaticEnv()        // 支持环境变量覆盖

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, reading environment variables from system")
	}

	// 设置默认值
	viper.SetDefault("APP_ENV", "dev")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGO_DB", "testdb")

	Cfg = &Config{
		AppEnv:           viper.GetString("APP_ENV"),
		AppPort:          viper.GetString("APP_PORT"),
		MongoURI:         viper.GetString("MONGO_URI"),
		MongoDB:          viper.GetString("MONGO_DB"),
		RedisAddr:        viper.GetString("REDIS_ADDR"),
		RedisPassword:    viper.GetString("REDIS_PASSWORD"),
		RedisDB:          viper.GetInt("REDIS_DB"),
		JWTAccessSecret:  viper.GetString("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
		JWTAccessExpire:  viper.GetInt("JWT_ACCESS_EXPIRE"),
		JWTRefreshExpire: viper.GetInt("JWT_REFRESH_EXPIRE"),
	}
}
