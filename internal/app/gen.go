package app

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/wzzfarewell/go-starter/internal/model"
	"github.com/wzzfarewell/go-starter/internal/util"
	"os"
	"strings"
	"text/template"
)

type Generator struct {
	ProjectPath string
	ModuleName  string
}

func NewGenerator(projectPath string, moduleName string) *Generator {
	return &Generator{
		ProjectPath: projectPath,
		ModuleName:  moduleName,
	}
}

func (g *Generator) Gen(configPath string) error {
	cfg, err := util.ReadConfigFromFile[model.Config](configPath)
	if err != nil {
		return errors.Wrap(err, "read model config from file failed")
	}
	appPath := fmt.Sprintf("%s/app", g.ProjectPath)
	if err := os.MkdirAll(appPath, 0755); err != nil {
		return errors.Wrap(err, "create app directory failed")
	}
	if err := g.genConfigFile(); err != nil {
		return errors.Wrap(err, "gen app config file failed")
	}
	if err := g.genRepositoryFile(cfg.Tables); err != nil {
		return errors.Wrap(err, "generate repository failed")
	}
	if err := g.genServiceFile(cfg.Tables); err != nil {
		return errors.Wrap(err, "generate service failed")
	}
	if err := g.genAppFile(); err != nil {
		return errors.Wrap(err, "generate controller failed")
	}
	return nil
}

func (g *Generator) genConfigFile() error {
	t := strings.ReplaceAll(ConfigTmpl, `'`, "`")
	tmpl, err := template.New("app config").Parse(t)
	if err != nil {
		return errors.Wrap(err, "parse app config template failed")
	}
	configFilePath := fmt.Sprintf("%s/app/config.go", g.ProjectPath)
	f, err := os.Create(configFilePath)
	if err != nil {
		return errors.Wrap(err, "create app config file failed")
	}
	if err := tmpl.Execute(f, nil); err != nil {
		return errors.Wrap(err, "execute app config template failed")
	}
	if err := f.Close(); err != nil {
		return errors.Wrap(err, "close app config file failed")
	}
	return nil
}

type RepositoryTmplText struct {
	ImportPkgPaths []string
	Repositories   []Repository
}

type Repository struct {
	PkgName string
	Name    string
}

func (g *Generator) genRepositoryFile(tables []model.Table) error {
	repos := make([]Repository, 0, len(tables))
	for _, v := range tables {
		repos = append(repos, Repository{
			PkgName: v.PackageName,
			Name:    v.StructName,
		})
	}
	text := RepositoryTmplText{
		ImportPkgPaths: []string{
			`"github.com/pkg/errors"`,
			`"github.com/wzzfarewell/go-mod/infrastructure/database"`,
			`"gorm.io/gorm"`,
		},
		Repositories: repos,
	}
	pkgNames := lo.Map(repos, func(t Repository, _ int) string {
		return t.PkgName
	})
	pkgNames = lo.Uniq(pkgNames)
	for _, v := range pkgNames {
		text.ImportPkgPaths = append(text.ImportPkgPaths, fmt.Sprintf(`"%s/domain/%s"`, g.ModuleName, v))
		text.ImportPkgPaths = append(text.ImportPkgPaths, fmt.Sprintf(`%sRepo "%s/repository/%s"`, v, g.ModuleName, v))
	}
	repoFilePath := fmt.Sprintf("%s/app/repository.go", g.ProjectPath)
	f, err := os.Create(repoFilePath)
	if err != nil {
		return errors.Wrap(err, "create app repository file failed")
	}
	tmpl, err := template.New("app repository").Parse(RepositoryTmpl)
	if err != nil {
		return errors.Wrap(err, "parse app repository template failed")
	}
	if err := tmpl.Execute(f, text); err != nil {
		return errors.Wrap(err, "execute app repository template failed")
	}
	if err := f.Close(); err != nil {
		return errors.Wrap(err, "close app repository file failed")
	}
	return nil
}

type ServiceTmplText struct {
	ImportPkgPaths []string
	Services       []Service
}

type Service struct {
	PkgName string
	Name    string
}

func (g *Generator) genServiceFile(tables []model.Table) error {
	services := make([]Service, 0, len(tables))
	for _, v := range tables {
		services = append(services, Service{
			PkgName: v.PackageName,
			Name:    v.StructName,
		})
	}
	text := ServiceTmplText{
		ImportPkgPaths: []string{},
		Services:       services,
	}
	pkgNames := lo.Map(services, func(t Service, _ int) string {
		return t.PkgName
	})
	pkgNames = lo.Uniq(pkgNames)
	for _, v := range pkgNames {
		text.ImportPkgPaths = append(text.ImportPkgPaths, fmt.Sprintf(`"%s/domain/%s"`, g.ModuleName, v))
		text.ImportPkgPaths = append(text.ImportPkgPaths, fmt.Sprintf(`%sSvc"%s/service/%s"`, v, g.ModuleName, v))
	}
	serviceFilePath := fmt.Sprintf("%s/app/service.go", g.ProjectPath)
	f, err := os.Create(serviceFilePath)
	if err != nil {
		return errors.Wrap(err, "create app service file failed")
	}
	tmpl, err := template.New("app service").Parse(ServiceTmpl)
	if err != nil {
		return errors.Wrap(err, "parse app service template failed")
	}
	if err := tmpl.Execute(f, text); err != nil {
		return errors.Wrap(err, "execute app service template failed")
	}
	if err := f.Close(); err != nil {
		return errors.Wrap(err, "close app service file failed")
	}
	return nil
}

type AppTmplText struct {
	ImportPkgPaths []string
}

func (g *Generator) genAppFile() error {
	text := AppTmplText{
		ImportPkgPaths: []string{
			`"github.com/gin-gonic/gin"`,
			`"github.com/pkg/errors"`,
			`"net/http"`,
			`"time"`,
			`"fmt"`,
		},
	}
	text.ImportPkgPaths = append(text.ImportPkgPaths, fmt.Sprintf(`"%s/delivery/http/middleware"`, g.ModuleName))
	appFilePath := fmt.Sprintf("%s/app/app.go", g.ProjectPath)
	f, err := os.Create(appFilePath)
	if err != nil {
		return errors.Wrap(err, "create app file failed")
	}
	tmpl, err := template.New("app").Parse(AppTmpl)
	if err != nil {
		return errors.Wrap(err, "parse app template failed")
	}
	if err := tmpl.Execute(f, text); err != nil {
		return errors.Wrap(err, "execute app template failed")
	}
	if err := f.Close(); err != nil {
		return errors.Wrap(err, "close app file failed")
	}
	return nil
}
