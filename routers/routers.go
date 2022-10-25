package routers

import (
	"github.com/jinzhu/gorm"
	"xcheck.info/telebot/pkg/api"
	"xcheck.info/telebot/pkg/repositories"
	"xcheck.info/telebot/pkg/services"
)

func InitUserAPI(db *gorm.DB) api.UserAPI {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userAPI := api.ProvideUserAPI(userService)
	return userAPI
}


func InitTransactionAPI(db *gorm.DB) api.TransactionAPI {
	transactionRepository := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepository)
	transactionAPI := api.ProvideTransactionAPI(transactionService)
	return transactionAPI
}
