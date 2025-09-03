package controllers

import (
	"net/http"
	"time"

	"github.com/amarjeetdev/ev-charging-app/db"
	"github.com/amarjeetdev/ev-charging-app/services"
	"github.com/gin-gonic/gin"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func HealthCheck(c *gin.Context) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	// Check PostgreSQL connection
	if db.DB != nil {
		sqlDB, err := db.DB.DB()
		if err != nil {
			status.Services["postgresql"] = "error: " + err.Error()
			status.Status = "unhealthy"
		} else {
			err = sqlDB.Ping()
			if err != nil {
				status.Services["postgresql"] = "error: " + err.Error()
				status.Status = "unhealthy"
			} else {
				status.Services["postgresql"] = "healthy"
			}
		}
	} else {
		status.Services["postgresql"] = "not connected"
		status.Status = "unhealthy"
	}

	// Check Redis connection (optional for development)
	if db.RedisClient != nil {
		_, err := db.RedisClient.Ping(db.Ctx).Result()
		if err != nil {
			status.Services["redis"] = "error: " + err.Error()
			status.Status = "unhealthy"
		} else {
			status.Services["redis"] = "healthy"
		}
	} else {
		status.Services["redis"] = "disabled (optional for development)"
		// Don't mark as unhealthy if Redis is optional
	}

	c.JSON(http.StatusOK, status)
}

func ReadinessCheck(c *gin.Context) {
	// Check if DB is connected (Redis is optional)
	if db.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database not initialized",
		})
		return
	}

	// Quick ping checks for database
	sqlDB, err := db.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database connection error",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database ping failed",
		})
		return
	}

	// Check Redis only if it's available
	if db.RedisClient != nil {
		if _, err := db.RedisClient.Ping(db.Ctx).Result(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"reason": "redis ping failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"timestamp": time.Now(),
	})
}

func GetCacheMetrics(c *gin.Context) {
	metrics, err := services.GetCacheMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func ResetCacheMetrics(c *gin.Context) {
	if err := services.ResetMetrics(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset cache metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cache metrics reset successfully"})
}
