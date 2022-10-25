package dtos

import (
	"github.com/jinzhu/gorm"
	"xcheck.info/telebot/pkg/models"
	_ "xcheck.info/telebot/pkg/models"
)

type UserRequest struct {
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

func ToUser(userDTO UserRequest) models.User {
	return models.User{
		ChatID:   userDTO.ChatID,
		UserID:   userDTO.UserID,
		UserName: userDTO.UserName,
	}
}

func ToUserDTO(user models.User) UserRequest {
	return UserRequest{
		ChatID:   user.ChatID,
		UserID:   user.UserID,
		UserName: user.UserName,
	}
}

func ToProductDTOs(products []models.User) []UserRequest {
	userdtos := make([]UserRequest, len(products))

	for i, itm := range products {
		userdtos[i] = ToUserDTO(itm)
	}

	return userdtos
}
