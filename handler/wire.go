package main

import (
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"xcheck.info/telebot/pkg/api"
	"xcheck.info/telebot/pkg/repositories"
	"xcheck.info/telebot/pkg/services"
)

func initProductAPI(db *gorm.DB) api.UserAPI {
	wire.Build(repositories.NewUserRepository, services.NewUserService, api.UserAPI{})

	return api.UserAPI{}
}
