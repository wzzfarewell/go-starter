package transport

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"text/template"
)

func (g *Generator) GenComponent() error {
	componentPath := g.ProjectPath + "/delivery/http/component"
	if err := os.MkdirAll(componentPath, 0755); err != nil {
		return errors.Wrap(err, "create component dir failed")
	}
	if err := g.genContextComponent(componentPath); err != nil {
		return errors.Wrap(err, "generate context component failed")
	}
	return nil
}

func (g *Generator) genContextComponent(componentPath string) error {
	contextComponentPath := componentPath + "/context.go"
	log.Printf("generate context component: %s", contextComponentPath)
	f, err := os.Create(contextComponentPath)
	if err != nil {
		return errors.Wrap(err, "create context component file failed")
	}
	tmpl, err := template.New("context component").Parse(ContextComponentTmpl)
	if err != nil {
		return errors.Wrap(err, "parse context component template failed")
	}
	if err := tmpl.Execute(f, nil); err != nil {
		return errors.Wrap(err, "execute context component template failed")
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
