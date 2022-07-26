package app

import (
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-mod/infrastructure/database"
	"github.com/wzzfarewell/go-starter-example/domain/user"
	userRepo "github.com/wzzfarewell/go-starter-example/repository/user"
	"gorm.io/gorm"
)

var (
	repositories Repositories
)

type Repositories struct {
	UserRepository     user.UserRepository
	UserInfoRepository user.UserInfoRepository
}

func initRepositories() error {
	db, err := initDB()
	if err != nil {
		return errors.Wrap(err, "init database failed")
	}
	repositories = Repositories{
		UserRepository:     userRepo.NewUserRepository(db),
		UserInfoRepository: userRepo.NewUserInfoRepository(db),
	}
	return nil
}

func initDB() (*gorm.DB, error) {
	return database.GormDB(&database.MySQLConfig{
		Host:     configs.DBConfig.Host,
		Port:     configs.DBConfig.Port,
		User:     configs.DBConfig.User,
		Password: configs.DBConfig.Password,
		DBName:   configs.DBConfig.DBName,
		LogLevel: configs.DBConfig.LogLevel,
	})
}
