package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
	Status     string `json:"status"`
	DeviceName string `json:"device_name"`
	ChatID     int    `json:"chat_id"`
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
}
