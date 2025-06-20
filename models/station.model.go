package models

import "gorm.io/gorm"

type Station struct {
	gorm.Model
	Name      string  `gorm:"not null" json:"name"`
	Address   string  `gorm:"not null" json:"address"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Bookings []Booking `gorm:"foreignKey:StationID" json:"bookings"`
}
