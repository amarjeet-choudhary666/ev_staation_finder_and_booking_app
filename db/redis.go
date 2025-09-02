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
		log.Fatal("❌ Failed to connect to Redis: ", err)
	}

	log.Println("✅ Redis connection established")
}

// Helper function to set cache with expiration
func SetCache(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(Ctx, key, value, expiration).Err()
}

// Helper function to get cache
func GetCache(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

// Helper function to delete cache
func DeleteCache(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}
