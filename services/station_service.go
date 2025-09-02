package services

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/amarjeetdev/ev-charging-app/db"
	"github.com/amarjeetdev/ev-charging-app/models"
)

const (
	StationCacheKey = "stations:all"
	StationCacheTTL = 10 * time.Minute
)

func GetAllStations() ([]models.Station, error) {
	// Try to get from cache first
	cachedStations, err := db.GetCache(StationCacheKey)
	if err == nil {
		var stations []models.Station
		if err := json.Unmarshal([]byte(cachedStations), &stations); err == nil {
			log.Println("✅ Retrieved stations from cache")
			RecordCacheHit()
			return stations, nil
		}
	}
	RecordCacheMiss()

	// If not in cache or error, get from database
	var stations []models.Station
	if err := db.DB.Find(&stations).Error; err != nil {
		return nil, err
	}

	// Cache the result
	if stationsJSON, err := json.Marshal(stations); err == nil {
		if err := db.SetCache(StationCacheKey, string(stationsJSON), StationCacheTTL); err != nil {
			log.Printf("⚠️ Failed to cache stations: %v", err)
		}
	}

	log.Println("✅ Retrieved stations from database")
	return stations, nil
}

func GetStationByID(id uint) (*models.Station, error) {
	cacheKey := "station:" + strconv.Itoa(int(id))

	// Try to get from cache first
	cachedStation, err := db.GetCache(cacheKey)
	if err == nil {
		var station models.Station
		if err := json.Unmarshal([]byte(cachedStation), &station); err == nil {
			log.Println("✅ Retrieved station from cache")
			RecordCacheHit()
			return &station, nil
		}
	}
	RecordCacheMiss()

	// If not in cache, get from database
	var station models.Station
	if err := db.DB.First(&station, id).Error; err != nil {
		return nil, err
	}

	// Cache the result
	if stationJSON, err := json.Marshal(station); err == nil {
		if err := db.SetCache(cacheKey, string(stationJSON), StationCacheTTL); err != nil {
			log.Printf("⚠️ Failed to cache station: %v", err)
		}
	}

	log.Println("✅ Retrieved station from database")
	return &station, nil
}

func CreateStation(station *models.Station) error {
	// Create in database
	if err := db.DB.Create(station).Error; err != nil {
		return err
	}

	// Invalidate cache
	if err := db.DeleteCache(StationCacheKey); err != nil {
		log.Printf("⚠️ Failed to invalidate stations cache: %v", err)
	}

	log.Println("✅ Created station and invalidated cache")
	return nil
}

func InvalidateStationCache(id uint) {
	cacheKey := "station:" + strconv.Itoa(int(id))
	if err := db.DeleteCache(cacheKey); err != nil {
		log.Printf("⚠️ Failed to invalidate station cache: %v", err)
	}
}
