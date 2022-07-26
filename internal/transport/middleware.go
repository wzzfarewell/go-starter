package transport

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"text/template"
)

func (g *Generator) GenMiddleware() error {
	middlewarePath := g.ProjectPath + "/delivery/http/middleware"
	if err := os.MkdirAll(middlewarePath, 0755); err != nil {
		return errors.Wrap(err, "create middleware dir failed")
	}
	if err := g.genErrorHandlerMiddleware(middlewarePath); err != nil {
		return errors.Wrap(err, "generate error handler middleware failed")
	}
	return nil
}

func (g *Generator) genErrorHandlerMiddleware(middlewarePath string) error {
	errorHandlerMiddlewarePath := middlewarePath + "/error.go"
	log.Printf("generate error handler middleware: %s", errorHandlerMiddlewarePath)
	f, err := os.Create(errorHandlerMiddlewarePath)
	if err != nil {
		return errors.Wrap(err, "create error handler middleware file failed")
	}
	tmpl, err := template.New("error handler middleware").Parse(ErrorHandlerMiddlewareTmpl)
	if err != nil {
		return errors.Wrap(err, "parse error handler middleware template failed")
	}
	if err := tmpl.Execute(f, nil); err != nil {
		return errors.Wrap(err, "execute error handler middleware template failed")
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
