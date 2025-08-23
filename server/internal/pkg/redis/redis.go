package redis

import (
	"context"
	"log"
	"time"

	"i18n-flow/internal/pkg/config"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.RedisAddr,
		Password: config.Cfg.RedisPassword,
		DB:       config.Cfg.RedisDB,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	log.Println("✅ Redis connected")
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func Del(key string) error {
	return Client.Del(Ctx).Err()
}
