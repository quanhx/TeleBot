package repositories

import (
	"github.com/jinzhu/gorm"
	"xcheck.info/telebot/pkg/models"
)

type UserRepository interface {
	FindByID(id uint) *models.User
	FindByUserID(userID uint) *models.User
	CreateUser(user *models.User) error
}

type userRepository struct {
	orm *gorm.DB
}

func NewUserRepository(orm *gorm.DB) UserRepository {
	return &userRepository{
		orm: orm,
	}
}

func (r *userRepository) FindByID(id uint) *models.User {
	var user models.User
	r.orm.First(&user, id)

	return &user
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.orm.Create(&user).Error
}

func (r *userRepository) FindByUserID(userID uint) *models.User {
	var user models.User
	r.orm.Model(&models.User{}).Where("deleted_at IS NULL").Where("user_id = ?", userID).Find(&user)

	return &user
}
