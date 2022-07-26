package api

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-starter/internal/model"
	"os"
	"text/template"
)

func genConfigFile(configPath, projectPath string) error {
	dbCfg, err := model.ReadDBConfig(configPath)
	if err != nil {
		return errors.Wrap(err, "read db config from file failed")
	}
	filePath := fmt.Sprintf("%s/application.toml", projectPath)
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "create application.toml failed")
	}
	tmpl, err := template.New("config").Parse(configTmpl)
	if err != nil {
		return errors.Wrap(err, "parse config template failed")
	}
	if err := tmpl.Execute(f, dbCfg); err != nil {
		return errors.Wrap(err, "execute config template failed")
	}
	if err := f.Close(); err != nil {
		return errors.Wrap(err, "close config file failed")
	}
	return nil
}

const configTmpl = `# This is a TOML config file.

[db]
db-name = "{{ .DBName }}"
host = "{{ .Host }}"
port = {{ .Port }}
user = "{{ .User }}"
password = "{{ .Password }}"

`
