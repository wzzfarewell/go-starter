package model

import (
	"fmt"
	"github.com/wzzfarewell/go-starter/internal/util"
)

type Config struct {
	DBConfig DBConfig `toml:"db" yaml:"db"`
	Tables   []Table  `toml:"tables" yaml:"tables"`
}

type DBConfig struct {
	DBName   string `toml:"db-name" yaml:"db-name"`
	Host     string `toml:"host" yaml:"host" default:"localhost"`
	Port     int    `toml:"port" yaml:"port" default:"3306"`
	User     string `toml:"user" yaml:"user"`
	Password string `toml:"password" yaml:"password"`
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName)
}

type Table struct {
	Name        string `toml:"name" yaml:"name"`
	StructName  string `toml:"struct-name" yaml:"struct-name"`
	PackageName string `toml:"package-name" yaml:"package-name"`
}

func ReadDBConfig(path string) (*DBConfig, error) {
	config, err := util.ReadConfigFromFile[Config](path)
	if err != nil {
		return nil, err
	}
	return &config.DBConfig, nil
}
