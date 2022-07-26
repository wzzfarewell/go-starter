package main

import (
	"github.com/urfave/cli/v2"
	"github.com/wzzfarewell/go-starter/api"
	"log"
	"os"
)

func main() {
	var projectPath, moduleName, configPath string
	a := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "projectPath",
				Aliases:     []string{"p"},
				Value:       "./go-starter-example",
				Usage:       "project path",
				Destination: &projectPath,
			},
			&cli.StringFlag{
				Name:        "moduleName",
				Aliases:     []string{"m"},
				Value:       "github.com/wzzfarewell/go-starter-example",
				Usage:       "module name",
				Destination: &moduleName,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "application.toml",
				Usage:       "generator configuration file",
				Destination: &configPath,
			},
		},
		Action: func(*cli.Context) error {
			return run(configPath)
		},
	}
	if err := a.Run(os.Args); err != nil {
		log.Fatal("app run failed: ", err)
	}
}

func run(configPath string) error {
	return api.NewGenerator(configPath).Gen()
}
