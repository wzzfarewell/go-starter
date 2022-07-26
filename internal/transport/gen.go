package transport

import "github.com/pkg/errors"

type Generator struct {
	ProjectPath string
}

func NewGenerator(projectPath string) *Generator {
	return &Generator{
		ProjectPath: projectPath,
	}
}

func (g *Generator) Gen() error {
	if err := g.GenComponent(); err != nil {
		return errors.Wrap(err, "generate component failed")
	}
	if err := g.GenMiddleware(); err != nil {
		return errors.Wrap(err, "generate middleware failed")
	}
	return nil
}
