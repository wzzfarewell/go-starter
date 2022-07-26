package main

import (
	"github.com/urfave/cli/v2"
	"github.com/wzzfarewell/go-mod/infrastructure/logger"
	"github.com/wzzfarewell/go-starter-example/app"
	"go.uber.org/zap"
	"os"
)

func main() {
	var cfgPath string
	a := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "application.toml",
				Usage:       "Load configuration from file",
				Destination: &cfgPath,
			},
		},
		Action: func(*cli.Context) error {
			return app.Run(cfgPath)
		},
	}

	if err := a.Run(os.Args); err != nil {
		logger.Fatal("app run failed", zap.Error(err))
	}
}
