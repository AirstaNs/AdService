package app

import (
	"homework10/internal/adapters/repository/adrepo"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/service"
	"homework10/internal/util"
)

//go:generate go run github.com/vektra/mockery/v2@v2.25.0 --name=App --filename=mockApp.go --output ../mocks/appemocks
type App interface {
	service.UserService
	service.AdService
}

type AdsApp struct {
	service.UserService
	service.AdService
}

func NewApp(adRepo adrepo.AdRepository, userRepo userrepo.UserRepository, formatter util.DateTimeFormatter) App {
	userService := service.NewUserService(userRepo)
	adService := service.NewAdsService(adRepo, formatter)
	return &AdsApp{userService, adService}
}
