package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // default fallback
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // default DB

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Ping to test connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Println("⚠️  Redis connection failed (optional for development): ", err)
		log.Println("ℹ️  Running without Redis - rate limiting and caching disabled")
		RedisClient = nil // Set to nil to indicate no connection
		return
	}

	log.Println("✅ Redis connection established")
}

// Helper function to set cache with expiration
func SetCache(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil // No-op if Redis not available
	}
	return RedisClient.Set(Ctx, key, value, expiration).Err()
}

// Helper function to get cache
func GetCache(key string) (string, error) {
	if RedisClient == nil {
		return "", redis.Nil // Return not found if Redis not available
	}
	return RedisClient.Get(Ctx, key).Result()
}

// Helper function to delete cache
func DeleteCache(key string) error {
	if RedisClient == nil {
		return nil // No-op if Redis not available
	}
	return RedisClient.Del(Ctx, key).Err()
}
