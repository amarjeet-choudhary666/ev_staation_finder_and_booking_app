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
	BookingCacheTTL = 5 * time.Minute
)

func GetUserBookings(userID uint) ([]models.Booking, error) {
	cacheKey := "user_bookings:" + strconv.Itoa(int(userID))

	// Try to get from cache first
	cachedBookings, err := db.GetCache(cacheKey)
	if err == nil {
		var bookings []models.Booking
		if err := json.Unmarshal([]byte(cachedBookings), &bookings); err == nil {
			log.Println("✅ Retrieved user bookings from cache")
			RecordCacheHit()
			return bookings, nil
		}
	}
	RecordCacheMiss()

	// If not in cache, get from database
	var bookings []models.Booking
	if err := db.DB.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}

	// Cache the result
	if bookingsJSON, err := json.Marshal(bookings); err == nil {
		if err := db.SetCache(cacheKey, string(bookingsJSON), BookingCacheTTL); err != nil {
			log.Printf("⚠️ Failed to cache user bookings: %v", err)
		}
	}

	log.Println("✅ Retrieved user bookings from database")
	return bookings, nil
}

func CreateBooking(booking *models.Booking) error {
	// Create in database
	if err := db.DB.Create(booking).Error; err != nil {
		return err
	}

	// Invalidate user bookings cache
	cacheKey := "user_bookings:" + strconv.Itoa(int(booking.UserID))
	if err := db.DeleteCache(cacheKey); err != nil {
		log.Printf("⚠️ Failed to invalidate user bookings cache: %v", err)
	}

	log.Println("✅ Created booking and invalidated cache")
	return nil
}

func CheckBookingConflicts(stationID uint, startTime, endTime time.Time) ([]models.Booking, error) {
	var conflicts []models.Booking
	if err := db.DB.Where("station_id = ? AND status != ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?))",
		stationID, "cancelled",
		endTime, startTime,
		endTime, startTime).
		Find(&conflicts).Error; err != nil {
		return nil, err
	}

	return conflicts, nil
}

func InvalidateUserBookingsCache(userID uint) {
	cacheKey := "user_bookings:" + strconv.Itoa(int(userID))
	if err := db.DeleteCache(cacheKey); err != nil {
		log.Printf("⚠️ Failed to invalidate user bookings cache: %v", err)
	}
}
