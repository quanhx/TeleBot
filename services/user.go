package services

import (
	"xcheck.info/telebot/pkg/dtos"
	"xcheck.info/telebot/pkg/models"
	"xcheck.info/telebot/pkg/repositories"
)

type UserService interface {
	CreateUser(req dtos.UserRequest) error
	FindByID(id uint) *models.User
	FindByUserID(userID uint) *models.User
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}

}

func (s *userService) CreateUser(req dtos.UserRequest) error {
	var user = models.User{}
	user = dtos.ToUser(req)
	return s.userRepo.CreateUser(&user)
}

func (s *userService) FindByID(id uint) *models.User {
	var response = s.userRepo.FindByID(id)
	if response == nil {
		return nil
	}
	return response
}

func (s *userService) FindByUserID(userID uint) *models.User {
	var response = s.userRepo.FindByUserID(userID)
	if response == nil {
		return nil
	}
	return response
}
