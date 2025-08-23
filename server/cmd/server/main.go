package main

import (
	"fmt"
	"log"

	"i18n-flow/internal/pkg/config"
	"i18n-flow/internal/pkg/logger"
	"i18n-flow/internal/pkg/mongo"
	"i18n-flow/internal/pkg/redis"
	"i18n-flow/internal/pkg/validator"
	"i18n-flow/internal/router"
)

func main() {
	// ===============================
	// 1️⃣ 初始化配置
	// ===============================
	config.Init()
	cfg := config.Cfg

	// ===============================
	// 2️⃣ 初始化日志
	// ===============================
	logger.Init(cfg.AppEnv)
	defer logger.Sync()
	logger.Info("Starting service...")

	// ===============================
	// 3️⃣ 初始化 Validator
	// ===============================
	validator.Init()
	logger.Info("Validator initialized")

	// ===============================
	// 4️⃣ 初始化 MongoDB
	// ===============================
	client := mongo.Init(cfg.MongoURI)
	if client == nil {
		log.Fatal("MongoDB init failed")
	}
	// 确保索引存在
	mongo.EnsureIndexes()
	logger.Info("MongoDB connected")

	// ===============================
	// 5️⃣ 初始化 Redis
	// ===============================
	redis.Init()
	logger.Info("Redis connected")

	// ===============================
	// 6️⃣ 初始化路由
	// ===============================
	r := router.SetupRouter()
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Info("Starting Gin server", logger.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Error("Server run failed", logger.String("error", err.Error()))
		log.Fatal(err)
	}
}
