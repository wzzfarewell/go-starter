package app

const ConfigTmpl = `package app

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
	DBConfig DBConfig	'toml:"db" yaml:"db"'
}

type DBConfig struct {
	Host     string 'toml:"host" yaml:"host" default:"localhost"'
	Port     int    'toml:"port" yaml:"port" default:"3306"'
	User     string 'toml:"user" yaml:"user"'
	Password string 'toml:"password" yaml:"password"'
	DBName   string 'toml:"db-name" yaml:"db-name"'
	LogLevel string 'toml:"log-level" yaml:"log-level" default:"silent"' // silent, error, warn, info
}

type LogConfigProperties struct {
	LogConfig LogConfig	'toml:"log" yaml:"log"'
}

type LogConfig struct {
	Level       string 'toml:"level" yaml:"level" default:"info"'
	Encoding    string 'toml:"encoding" yaml:"encoding" default:"json"'
	Development bool   'toml:"development" yaml:"development" default:"false"'
	FileName    string 'toml:"file-name" yaml:"file-name" default:"info.log"'
	MaxSize     int    'toml:"max-size" yaml:"max-size" default:"128"'
	MaxBackups  int    'toml:"max-backups" yaml:"max-backups" default:"30"'
	MaxAge      int    'toml:"max-age" yaml:"max-age" default:"7"'
	Compress    bool   'toml:"compress" yaml:"compress" default:"false"'
}

type HTTPConfigProperties struct {
	HTTPConfig HTTPConfig	'toml:"http" yaml:"http"'
}	

type HTTPConfig struct {
	Port         int    'toml:"port" yaml:"port" default:"8080"'
	WriteTimeout int    'toml:"write-timeout" yaml:"write-timeout" default:"10"'
	ReadTimeout  int    'toml:"read-timeout" yaml:"read-timeout" default:"10"'
	Mode         string 'toml:"mode" yaml:"mode" default:"debug"'
}

`

const RepositoryTmpl = `package app

import (
	{{- range .ImportPkgPaths}}
	{{ printf "%s" . }}
	{{- end}}
)

var (
	repositories Repositories
)

type Repositories struct {
	{{- range $repo := .Repositories }}
	{{ $repo.Name }}Repository {{ $repo.PkgName }}.{{ $repo.Name }}Repository
	{{- end }}
}

func initRepositories() error {
	db, err := initDB()
	if err != nil {
		return errors.Wrap(err, "init database failed")
	}
	repositories = Repositories{
		{{- range $repo := .Repositories }}
		{{ $repo.Name }}Repository: {{ $repo.PkgName }}Repo.New{{ $repo.Name }}Repository(db),
		{{- end }}
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

`

const ServiceTmpl = `package app

import (
	{{- range .ImportPkgPaths}}
	{{ printf "%s" . }}
	{{- end}}
)

var (
	services Services
)

type Services struct {
	{{- range $svc := .Services }}
	{{ $svc.Name }}Service {{ $svc.PkgName }}.{{ $svc.Name }}Service
	{{- end }}
}

func initServices() error {
	services = Services{
		{{- range $svc := .Services }}
		{{ $svc.Name }}Service: {{ $svc.PkgName }}Svc.New{{ $svc.Name }}Service(repositories.{{ $svc.Name }}Repository),
		{{- end }}
	}
	return nil
}

`

const AppTmpl = `package app

import (
	{{- range .ImportPkgPaths}}
	{{ printf "%s" . }}
	{{- end}}
)

func Run(path string) error {
	if err := initConfigs(path); err != nil {
		return errors.Wrap(err, "init config failed")
	}
	if err := initRepositories(); err != nil {
		return errors.Wrap(err, "init repository failed")
	}
	if err := initServices(); err != nil {
		return errors.Wrap(err, "init service failed")
	}
	return serveHTTP()
}

func serveHTTP() error {
	e := gin.New()
	_ = e.SetTrustedProxies(nil)
	e.Use(gin.Recovery())
	e.NoRoute(middleware.NoRouteHandler())
	e.Use(middleware.ErrorHandler())
	cfg := configs.HTTPConfig
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      e,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
	}
	return s.ListenAndServe()
}

`
