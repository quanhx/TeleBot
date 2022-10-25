package services

import (
	"go.uber.org/dig"
	"xcheck.info/telebot/pkg/database"
	"xcheck.info/telebot/pkg/repositories"
)

var serviceContainer *dig.Container

func InitServices()  {
	container := dig.New()

	_ = container.Provide(database.InitDb())
	_ = container.Provide(repositories.NewUserRepository)
	_ = container.Provide(NewUserService)

	serviceContainer = container

}

func GetServiceContainer() *dig.Container {
	return serviceContainer
}
