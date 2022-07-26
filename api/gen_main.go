package api

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"text/template"
)

func genMainFile(projectPath, moduleName string) error {
	text := MainTmplText{
		ImportPkgPaths: []string{
			`"github.com/urfave/cli/v2"`,
			`"github.com/wzzfarewell/go-mod/infrastructure/logger"`,
			`"go.uber.org/zap"`,
			`"os"`,
		},
	}
	text.ImportPkgPaths = append(text.ImportPkgPaths, fmt.Sprintf(`"%s/app"`, moduleName))
	mainFilePath := fmt.Sprintf(`%s/main.go`, projectPath)
	f, err := os.Create(mainFilePath)
	if err != nil {
		return errors.Wrap(err, "create main file failed")
	}
	tmpl, err := template.New("main").Parse(MainTmpl)
	if err != nil {
		return errors.Wrap(err, "parse main template failed")
	}
	if err := tmpl.Execute(f, text); err != nil {
		return errors.Wrap(err, "execute main template failed")
	}
	if err := f.Close(); err != nil {
		return errors.Wrap(err, "close main file failed")
	}
	return nil
}

type MainTmplText struct {
	ImportPkgPaths []string
}

const MainTmpl = `package main

import (
	{{- range .ImportPkgPaths }}
	{{ printf "%s" . }}
	{{- end}}
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

`
