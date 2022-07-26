package model

import (
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/wzzfarewell/go-starter/internal/util"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type Generator struct {
	projectPath string
	moduleName  string
}

func NewGenerator(projectPath, moduleName string) *Generator {
	return &Generator{
		projectPath: projectPath,
		moduleName:  moduleName,
	}
}

func (g *Generator) Gen(configPath string) error {
	cfg, err := util.ReadConfigFromFile[Config](configPath)
	if err != nil {
		return errors.Wrap(err, "read model config from file failed")
	}
	tables := lo.GroupBy(cfg.Tables, func(t Table) string {
		return t.PackageName
	})
	if err := g.GenModel(cfg.DBConfig, tables); err != nil {
		return errors.Wrap(err, "generate model failed")
	}
	if err := g.GenDomain(tables); err != nil {
		return errors.Wrap(err, "generate domain failed")
	}
	if err := g.GenGoMod(); err != nil {
		return errors.Wrap(err, "generate go mod failed")
	}
	return nil
}

func (g *Generator) GenGoMod() error {
	goModPath := g.projectPath + "/go.mod"
	f, err := os.Create(goModPath)
	if err != nil {
		return errors.Wrap(err, "create go mod file failed")
	}
	tmpl, err := template.New("go mod").Parse(GoModTmpl)
	if err != nil {
		return errors.Wrap(err, "parse go mod template failed")
	}
	if err := tmpl.Execute(f, g.moduleName); err != nil {
		return errors.Wrap(err, "execute go mod template failed")
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) GenModel(dbConfig DBConfig, tables map[string][]Table) error {
	db, err := gorm.Open(mysql.Open(dbConfig.DSN()), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "open database connection failed")
	}
	for k, t := range tables {
		// see more: https://github.com/go-gorm/gen/blob/master/README.ZH_CN.md#generate
		generator := gen.NewGenerator(gen.Config{
			OutPath:       fmt.Sprintf("%s/repository/%s/query", g.projectPath, k),
			Mode:          gen.WithoutContext,
			FieldNullable: true,
		})
		generator.UseDB(db)
		for _, v := range t {
			generator.ApplyBasic(generator.GenerateModelAs(v.Name, v.StructName))
		}
		// execute the action of code generation
		generator.Execute()
		for _, v := range t {
			path := fmt.Sprintf("%s/repository/%s/query/%s.gen.go", g.projectPath, v.PackageName, v.Name)
			data, err := os.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "read generated model file failed")
			}
			text := string(data)
			// replace the package name
			modelPkgName := fmt.Sprintf(`"%s/repository/%s/model"`, g.moduleName, v.PackageName)
			text = regexp.MustCompile(`"(.*?)/model"`).ReplaceAllString(text, modelPkgName)
			if err := os.WriteFile(path, []byte(text), 0644); err != nil {
				return errors.Wrap(err, "write generated model file failed")
			}
		}
	}
	return nil
}

func (g *Generator) GenDomain(tables map[string][]Table) error {
	repoPath := g.projectPath + "/repository"
	for k, v := range tables {
		modelDirPath := fmt.Sprintf("%s/%s/model", repoPath, k)
		domainDirPath := fmt.Sprintf("%s/domain/%s", g.projectPath, k)
		if err := os.MkdirAll(domainDirPath, 0755); err != nil {
			return errors.Wrap(err, "create domain dir failed")
		}
		serviceDirPath := fmt.Sprintf("%s/service/%s", g.projectPath, k)
		if err := os.MkdirAll(serviceDirPath, 0755); err != nil {
			return errors.Wrap(err, "create service dir failed")
		}
		for _, t := range v {
			if err := g.genDomainFile(modelDirPath, domainDirPath, k, t.StructName, t.Name); err != nil {
				return errors.Wrap(err, "gen domain file failed")
			}
			if err := g.genRepositoryFile(k, t.StructName); err != nil {
				return errors.Wrap(err, "gen repository file failed")
			}
			if err := g.GenServiceFile(k, t.StructName); err != nil {
				return errors.Wrap(err, "gen service file failed")
			}
		}

	}
	return nil
}

func (g *Generator) genDomainFile(modelDirPath, domainDirPath, pkgName, structName, tableName string) error {
	modelFilePath := modelDirPath + "/" + tableName + ".gen.go"
	data, err := os.ReadFile(modelFilePath)
	if err != nil {
		log.Printf("read model file %s failed, error: %v", modelFilePath, err)
		return nil
	}
	content := string(data)
	domainFileName := stringy.New(structName).SnakeCase().ToLower()
	domainFilePath := fmt.Sprintf("%s/%s.go", domainDirPath, domainFileName)
	log.Printf("generate domain file: %s", domainFilePath)
	f, err := os.Create(domainFilePath)
	if err != nil {
		return errors.Wrap(err, "create domain file failed")
	}
	text := DomainTmplText{
		PkgName:        pkgName,
		ImportPkgPaths: []string{"context"},
		StructName:     structName,
	}
	start := strings.Index(content, "type")
	end := strings.Index(content, "}")
	if start == -1 || end == -1 {
		log.Printf("parse model file %s failed, struct definition not found", modelFilePath)
		return nil
	}
	typeStructStr := content[start : end+1]
	r := regexp.MustCompile(`gorm:"(.*?)" `)
	typeStructStr = r.ReplaceAllString(typeStructStr, "")
	text.Struct = typeStructStr
	if strings.Contains(typeStructStr, "time.Time") {
		text.ImportPkgPaths = append(text.ImportPkgPaths, "time")
	}
	tmpl, err := template.New("domain").Parse(DomainTmpl)
	if err != nil {
		return errors.Wrap(err, "parse domain template failed")
	}
	if err := tmpl.Execute(f, text); err != nil {
		return errors.Wrap(err, "execute domain template failed")
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) genRepositoryFile(domainPkgName, structName string) error {
	repoFileName := stringy.New(structName).SnakeCase().ToLower()
	repoFilePath := fmt.Sprintf("%s/repository/%s/%s.go", g.projectPath, domainPkgName, repoFileName)
	log.Printf("generate repository file: %s", repoFilePath)
	f, err := os.Create(repoFilePath)
	if err != nil {
		return errors.Wrap(err, "create repository file failed")
	}
	text := RepositoryTmplText{
		PkgName: domainPkgName,
		ImportPkgPaths: []string{
			"github.com/wzzfarewell/go-mod/infrastructure/database",
			fmt.Sprintf("%s/domain/%s", g.moduleName, domainPkgName),
			fmt.Sprintf("%s/repository/%s/model", g.moduleName, domainPkgName),
			fmt.Sprintf("%s/repository/%s/query", g.moduleName, domainPkgName),
			"gorm.io/gorm",
		},
		StructName:    structName,
		DomainPkgName: domainPkgName,
	}
	tmpl, err := template.New("repository").Parse(RepositoryTmpl)
	if err != nil {
		return errors.Wrap(err, "parse repository template failed")
	}
	if err := tmpl.Execute(f, text); err != nil {
		return errors.Wrap(err, "execute repository template failed")
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) GenServiceFile(pkgName, structName string) error {
	serviceFileName := stringy.New(structName).SnakeCase().ToLower()
	serviceFilePath := fmt.Sprintf("%s/service/%s/%s.go", g.projectPath, pkgName, serviceFileName)
	log.Printf("generate service file: %s", serviceFilePath)
	f, err := os.Create(serviceFilePath)
	if err != nil {
		return errors.Wrap(err, "create service file failed")
	}
	text := ServiceTmplText{
		PkgName: pkgName,
		ImportPkgPaths: []string{
			"context",
			fmt.Sprintf("%s/domain/%s", g.moduleName, pkgName),
			"github.com/wzzfarewell/go-mod/infrastructure/utils/copyutil",
			"github.com/pkg/errors",
		},
		StructName: structName,
	}
	tmpl, err := template.New("service").Parse(ServiceTmpl)
	if err != nil {
		return errors.Wrap(err, "parse service template failed")
	}
	if err := tmpl.Execute(f, text); err != nil {
		return errors.Wrap(err, "execute service template failed")
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
