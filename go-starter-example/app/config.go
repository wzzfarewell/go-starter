package app

import (
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-mod/infrastructure/utils/fileutil"
)

var (
	configs Configs
)

type Configs struct {
	DBConfig   DBConfig
	LogConfig  LogConfig
	HTTPConfig HTTPConfig
}

func initConfigs(path string) error {
	dbConfig, err := fileutil.ReadConfigFromFile[DBConfigProperties](path)
	if err != nil {
		return errors.Wrap(err, "read db config failed")
	}
	logConfig, err := fileutil.ReadConfigFromFile[LogConfigProperties](path)
	if err != nil {
		return errors.Wrap(err, "read log config failed")
	}
	httpConfig, err := fileutil.ReadConfigFromFile[HTTPConfigProperties](path)
	if err != nil {
		return errors.Wrap(err, "read http config failed")
	}
	configs = Configs{
		DBConfig:   dbConfig.DBConfig,
		LogConfig:  logConfig.LogConfig,
		HTTPConfig: httpConfig.HTTPConfig,
	}
	return nil
}

type DBConfigProperties struct {
	DBConfig DBConfig `toml:"db" yaml:"db"`
}

type DBConfig struct {
	Host     string `toml:"host" yaml:"host" default:"localhost"`
	Port     int    `toml:"port" yaml:"port" default:"3306"`
	User     string `toml:"user" yaml:"user"`
	Password string `toml:"password" yaml:"password"`
	DBName   string `toml:"db-name" yaml:"db-name"`
	LogLevel string `toml:"log-level" yaml:"log-level" default:"silent"` // silent, error, warn, info
}

type LogConfigProperties struct {
	LogConfig LogConfig `toml:"log" yaml:"log"`
}

type LogConfig struct {
	Level       string `toml:"level" yaml:"level" default:"info"`
	Encoding    string `toml:"encoding" yaml:"encoding" default:"json"`
	Development bool   `toml:"development" yaml:"development" default:"false"`
	FileName    string `toml:"file-name" yaml:"file-name" default:"info.log"`
	MaxSize     int    `toml:"max-size" yaml:"max-size" default:"128"`
	MaxBackups  int    `toml:"max-backups" yaml:"max-backups" default:"30"`
	MaxAge      int    `toml:"max-age" yaml:"max-age" default:"7"`
	Compress    bool   `toml:"compress" yaml:"compress" default:"false"`
}

type HTTPConfigProperties struct {
	HTTPConfig HTTPConfig `toml:"http" yaml:"http"`
}

type HTTPConfig struct {
	Port         int    `toml:"port" yaml:"port" default:"8080"`
	WriteTimeout int    `toml:"write-timeout" yaml:"write-timeout" default:"10"`
	ReadTimeout  int    `toml:"read-timeout" yaml:"read-timeout" default:"10"`
	Mode         string `toml:"mode" yaml:"mode" default:"debug"`
}
