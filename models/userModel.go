package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   string `gorm:"unique"`
	Password string
}
