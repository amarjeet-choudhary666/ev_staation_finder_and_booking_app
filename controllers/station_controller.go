package controllers

import (
	"net/http"

	"github.com/amarjeetdev/ev-charging-app/db"
	"github.com/amarjeetdev/ev-charging-app/models"
	"github.com/gin-gonic/gin"
)

func CreateStation(c *gin.Context) {
	var station models.Station

	if err := c.ShouldBindBodyWithJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("address = ?", station.Address).First(&station).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "station already exists"})
		return
	}

	if err := db.DB.Create(&station).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create station"})
		return
	}

	c.JSON(http.StatusCreated, station)
}

func GetAllStations(c *gin.Context) {
	var station []models.Station

	if err := db.DB.Find(&station).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}

	if err := db.DB.Find(&station).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stations"})
		return
	}

	c.JSON(http.StatusOK, station)
}

func GetStationDetails(c *gin.Context) {
	id := c.Param("id")
	var station models.Station

	if err := db.DB.First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	c.JSON(http.StatusOK, station)
}
