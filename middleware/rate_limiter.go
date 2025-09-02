package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/amarjeetdev/ev-charging-app/db"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	RateLimitWindow = 15 * time.Minute
	RateLimitMax    = 100 // requests per window
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP

		// Get current count
		count, err := db.RedisClient.Get(db.Ctx, key).Result()
		if err == redis.Nil {
			// First request, set count to 1
			err = db.RedisClient.Set(db.Ctx, key, 1, RateLimitWindow).Err()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
				return
			}
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			return
		} else {
			// Increment count
			currentCount, _ := strconv.Atoi(count)
			if currentCount >= RateLimitMax {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
					"retry_after": RateLimitWindow.Seconds(),
				})
				return
			}

			err = db.RedisClient.Incr(db.Ctx, key).Err()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
				return
			}
		}

		c.Next()
	}
}
