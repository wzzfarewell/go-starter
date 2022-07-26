package model

type DomainTmplText struct {
	PkgName        string
	ImportPkgPaths []string
	StructName     string
	Struct         string
}

const DomainTmpl = `package {{ .PkgName }}

import (
	{{- range .ImportPkgPaths}}
	{{ printf "%q" . }}
	{{- end}}
)

{{ .Struct }}

type {{ .StructName }}Repository interface {
	Create(entity ...*{{ .StructName }}) error
	Update(entity *{{ .StructName }}) error
	Delete(id ...int32) (int64, error)
	Get(id int32) (*{{ .StructName }}, error)
	FindAll() ([]*{{ .StructName }}, error)
}

type {{ .StructName }}Service interface {
	Create(ctx context.Context, entity ...*{{ .StructName }}) error
	Update(ctx context.Context, entity *{{ .StructName }}) (*{{ .StructName }}, error)
	Delete(ctx context.Context, id ...int32) error
	Get(ctx context.Context, id int32) (*{{ .StructName }}, error)
	FindAll(ctx context.Context) ([]*{{ .StructName }}, error)
}

`

type RepositoryTmplText struct {
	PkgName        string
	ImportPkgPaths []string
	DomainPkgName  string
	StructName     string
}

const RepositoryTmpl = `package {{ .PkgName }}

import (
	{{- range .ImportPkgPaths}}
	{{ printf "%q" . }}
	{{- end}}
)

type {{ .StructName }}Repository struct {
	q *query.Query
	database.ModelConverter[model.{{ .StructName }}, {{ .DomainPkgName }}.{{ .StructName }}]
}

func New{{ .StructName }}Repository(db *gorm.DB) {{ .DomainPkgName }}.{{ .StructName }}Repository {
	return &{{ .StructName }}Repository{
		q: query.Use(db),
	}
}

func (r *{{ .StructName }}Repository) Get(id int32) (*{{ .DomainPkgName }}.{{ .StructName }}, error) {
	m, err := r.q.{{ .StructName }}.Where(r.q.{{ .StructName }}.ID.Eq(id)).First()
	return r.ToDomainWithError(m, err)
}

func (r *{{ .StructName }}Repository) Create(entity ...*{{ .DomainPkgName }}.{{ .StructName }}) error {
	models, err := r.ToModels(entity)
	if err != nil {
		return err
	}
	return r.q.{{ .StructName }}.Create(models...)
}

func (r *{{ .StructName }}Repository) Update(entity *{{ .DomainPkgName }}.{{ .StructName }}) error {
	_, err := r.q.{{ .StructName }}.Updates(entity)
	return err
}

func (r *{{ .StructName }}Repository) Delete(id ...int32) (int64, error) {
	info, err := r.q.{{ .StructName }}.Where(r.q.{{ .StructName }}.ID.In(id...)).Delete()
	return info.RowsAffected, err
}

func (r *{{ .StructName }}Repository) FindAll() ([]*{{ .DomainPkgName }}.{{ .StructName }}, error) {
	models, err := r.q.{{ .StructName }}.Find()
	if err != nil {
		return nil, err
	}
	return r.ToDomains(models)
}

`

type ServiceTmplText struct {
	PkgName        string
	ImportPkgPaths []string
	StructName     string
}

const ServiceTmpl = `package {{ .PkgName }}

import (
	{{- range .ImportPkgPaths}}
	{{ printf "%q" . }}
	{{- end}}
)

type {{ .StructName }}Service struct {
	{{ .StructName }}Repo {{ .PkgName }}.{{ .StructName }}Repository
}

func New{{ .StructName }}Service({{ .StructName }}Repo {{ .PkgName }}.{{ .StructName }}Repository) {{ .PkgName }}.{{ .StructName }}Service {
	return &{{ .StructName }}Service{
		{{ .StructName }}Repo: {{ .StructName }}Repo,
	}
}

func (s *{{ .StructName }}Service) Get(ctx context.Context, id int32) (*{{ .PkgName }}.{{ .StructName }}, error) {
	return s.{{ .StructName }}Repo.Get(id)
}

func (s *{{ .StructName }}Service) Create(ctx context.Context, entity ...*{{ .PkgName }}.{{ .StructName }}) error {
	return s.{{ .StructName }}Repo.Create(entity...)
}

func (s *{{ .StructName }}Service) Update(ctx context.Context, entity *{{ .PkgName }}.{{ .StructName }}) (*{{ .PkgName }}.{{ .StructName }}, error) {
	m, err := s.{{ .StructName }}Repo.Get(entity.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get model from repository")
	}
	if m == nil {
		return nil, errors.New("model not found")
	}
	if err := copyutil.CopyTo(entity, m); err != nil {
		return nil, errors.Wrap(err, "failed to copy entity to model")
	}
	if err := s.{{ .StructName }}Repo.Update(m); err != nil {
		return nil, errors.Wrap(err, "failed to update model in repository")
	}
	return m, nil
}

func (s *{{ .StructName }}Service) Delete(ctx context.Context, id ...int32) error {
	_, err := s.{{ .StructName }}Repo.Delete(id...)
	return err
}

func (s *{{ .StructName }}Service) FindAll(ctx context.Context) ([]*{{ .PkgName }}.{{ .StructName }}, error) {
	return s.{{ .StructName }}Repo.FindAll()
}

`

const GoModTmpl = `module {{ . }}

go 1.18

`
