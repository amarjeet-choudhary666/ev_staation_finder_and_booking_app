package controllers

import (
	"net/http"
	"strconv"

	"github.com/amarjeetdev/ev-charging-app/db"
	"github.com/amarjeetdev/ev-charging-app/models"
	"github.com/amarjeetdev/ev-charging-app/services"
	"github.com/gin-gonic/gin"
)

func CreateStation(c *gin.Context) {
	var station models.Station

	if err := c.ShouldBindBodyWithJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if station already exists
	var existingStation models.Station
	if err := db.DB.Where("address = ?", station.Address).First(&existingStation).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "station already exists"})
		return
	}

	if err := services.CreateStation(&station); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create station"})
		return
	}

	c.JSON(http.StatusCreated, station)
}

func GetAllStations(c *gin.Context) {
	stations, err := services.GetAllStations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stations"})
		return
	}

	c.JSON(http.StatusOK, stations)
}

func GetStationDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid station id"})
		return
	}

	station, err := services.GetStationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	c.JSON(http.StatusOK, station)
}
