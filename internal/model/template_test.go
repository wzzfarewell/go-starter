package model

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"text/template"
)

func TestDomainTmpl(t *testing.T) {
	s := `
type User struct {
	ID        int32     
	Name      string    
	CreatedAt time.Time 
	UpdatedAt time.Time 
}
`
	a := DomainTmplText{
		PkgName:        "user",
		ImportPkgPaths: []string{"context"},
		StructName:     "User",
		Struct:         s,
	}
	if strings.Contains(s, "time.Time") {
		a.ImportPkgPaths = append(a.ImportPkgPaths, "time")
	}
	tmpl, err := template.New("").Parse(DomainTmpl)
	require.NoError(t, err)
	var b strings.Builder
	err = tmpl.Execute(&b, a)
	require.NoError(t, err)
	t.Log(b.String())
}
