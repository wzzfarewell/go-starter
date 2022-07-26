package app

import (
	"github.com/wzzfarewell/go-starter-example/domain/user"
	userSvc "github.com/wzzfarewell/go-starter-example/service/user"
)

var (
	services Services
)

type Services struct {
	UserService     user.UserService
	UserInfoService user.UserInfoService
}

func initServices() error {
	services = Services{
		UserService:     userSvc.NewUserService(repositories.UserRepository),
		UserInfoService: userSvc.NewUserInfoService(repositories.UserInfoRepository),
	}
	return nil
}
