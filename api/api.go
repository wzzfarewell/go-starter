package api

import (
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-starter/internal/app"
	"github.com/wzzfarewell/go-starter/internal/model"
	"github.com/wzzfarewell/go-starter/internal/transport"
	"github.com/wzzfarewell/go-starter/internal/util"
)

type Generator interface {
	Gen() error
}

type generator struct {
	configPath string
}

func NewGenerator(configPath string) Generator {
	return &generator{
		configPath: configPath,
	}
}

func (g *generator) Gen() error {
	cfg, err := util.ReadConfigFromFile[Config](g.configPath)
	if err != nil {
		return errors.Wrap(err, "read config from file failed")
	}
	modelGen := model.NewGenerator(cfg.ProjectPath, cfg.ModuleName)
	if err := modelGen.Gen(g.configPath); err != nil {
		return errors.Wrap(err, "generate model failed")
	}
	transportGen := transport.NewGenerator(cfg.ProjectPath)
	if err := transportGen.Gen(); err != nil {
		return errors.Wrap(err, "generate transport failed")
	}
	appGen := app.NewGenerator(cfg.ProjectPath, cfg.ModuleName)
	if err := appGen.Gen(g.configPath); err != nil {
		return errors.Wrap(err, "generate app failed")
	}
	if err := genMainFile(cfg.ProjectPath, cfg.ModuleName); err != nil {
		return errors.Wrap(err, "generate main file failed")
	}
	if err := genConfigFile(g.configPath, cfg.ProjectPath); err != nil {
		return errors.Wrap(err, "generate config file failed")
	}
	return nil
}

type Config struct {
	ProjectPath string `toml:"project-path" yaml:"project-path"`
	ModuleName  string `toml:"module-name" yaml:"module-name"`
}
