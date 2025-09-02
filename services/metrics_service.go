package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/amarjeetdev/ev-charging-app/db"
)

type CacheMetrics struct {
	TotalRequests    int64   `json:"total_requests"`
	CacheHits        int64   `json:"cache_hits"`
	CacheMisses      int64   `json:"cache_misses"`
	HitRate          float64 `json:"hit_rate"`
	LastUpdated      time.Time `json:"last_updated"`
}

const MetricsKey = "cache_metrics"

func RecordCacheHit() {
	updateMetrics(true)
}

func RecordCacheMiss() {
	updateMetrics(false)
}

func updateMetrics(isHit bool) {
	metrics, err := getMetrics()
	if err != nil {
		log.Printf("⚠️ Failed to get metrics: %v", err)
		return
	}

	metrics.TotalRequests++
	if isHit {
		metrics.CacheHits++
	} else {
		metrics.CacheMisses++
	}

	if metrics.TotalRequests > 0 {
		metrics.HitRate = float64(metrics.CacheHits) / float64(metrics.TotalRequests) * 100
	}

	metrics.LastUpdated = time.Now()

	if err := saveMetrics(metrics); err != nil {
		log.Printf("⚠️ Failed to save metrics: %v", err)
	}
}

func getMetrics() (*CacheMetrics, error) {
	data, err := db.GetCache(MetricsKey)
	if err != nil {
		// Return default metrics if not found
		return &CacheMetrics{
			TotalRequests: 0,
			CacheHits:     0,
			CacheMisses:   0,
			HitRate:       0,
			LastUpdated:   time.Now(),
		}, nil
	}

	var metrics CacheMetrics
	if err := json.Unmarshal([]byte(data), &metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}

func saveMetrics(metrics *CacheMetrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return db.SetCache(MetricsKey, string(data), 24*time.Hour) // Keep metrics for 24 hours
}

func GetCacheMetrics() (*CacheMetrics, error) {
	return getMetrics()
}

func ResetMetrics() error {
	metrics := &CacheMetrics{
		TotalRequests: 0,
		CacheHits:     0,
		CacheMisses:   0,
		HitRate:       0,
		LastUpdated:   time.Now(),
	}

	return saveMetrics(metrics)
}
