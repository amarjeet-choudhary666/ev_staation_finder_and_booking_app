package controllers

import (
	"net/http"
	"time"

	"github.com/amarjeetdev/ev-charging-app/db"
	"github.com/amarjeetdev/ev-charging-app/models"
	"github.com/gin-gonic/gin"
)

func CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !booking.StartTime.Before(booking.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
		return
	}

	var conflicts []models.Booking
	if err := db.DB.Where("station_id = ? AND status != ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?))",
		booking.StationID, "cancelled",
		booking.EndTime, booking.StartTime,
		booking.EndTime, booking.StartTime).
		Find(&conflicts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for conflicts"})
		return
	}

	if len(conflicts) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Time slot already booked"})
		return
	}

	if booking.StartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be in the future"})
		return
	}

	if err := db.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func GetUserBookings(c *gin.Context) {
	userId := c.Query("user_id")

	var bookings []models.Booking
	if err := db.DB.Where("user_id = ?", userId).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
