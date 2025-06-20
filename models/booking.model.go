package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	StationID uint      `gorm:"not null" json:"station_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`

	Station Station `gorm:"foreignKey:StationID" json:"station"`
}
