package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	BookingID uint   `json:"booking_id"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
}
