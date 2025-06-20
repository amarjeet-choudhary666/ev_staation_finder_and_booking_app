package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100)" json:"name"`
	Email        string `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	Password     string `gorm:"type:varchar(255)" json:"password"`
	Role         string `gorm:"type:varchar(20);default:user" json:"role"`
	RefreshToken string `gorm:"type:varchar(255)" json:"refresh_token"`
}
