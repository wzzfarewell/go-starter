package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/wzzfarewell/go-mod/infrastructure/logger"
	"regexp"
	"strings"
	"testing"
)

func TestReadConfigFromFile(t *testing.T) {
	type dbConfig struct {
		DBName   string `toml:"db-name" yaml:"db-name"`
		Host     string `toml:"host" yaml:"host" default:"localhost"`
		Port     int    `toml:"port" yaml:"port" default:"3306"`
		User     string `toml:"user" yaml:"user"`
		Password string `toml:"password" yaml:"password"`
	}
	type config struct {
		DBConfig dbConfig `toml:"db" yaml:"db"`
	}
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "success",
			path:    "../testdata/config.toml",
			wantErr: false,
		},
		{
			name:    "success with not exist file",
			path:    "./testdata/not-exist.toml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ReadConfigFromFile[config](tt.path)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				t.Logf("%+v", c)
			}
		})
	}
}

func TestRegex(t *testing.T) {
	str := `package model

import (
	"time"
)

const TableNameUser = "t_user"

// User mapped from table <t_user>
type User struct {
	ID        int32     gorm:"column:id;primaryKey;autoIncrement:true" json:"id"
	Name      string    
	CreatedAt time.Time 
	UpdatedAt time.Time 
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}`
	start := strings.Index(str, "type")
	end := strings.Index(str, "}")
	s := str[start : end+1]
	r := regexp.MustCompile(`gorm:"(.*?)"`)
	t.Logf("%+v", r.FindAllString(s, -1))
	s = r.ReplaceAllString(s, "")
	t.Log(s)
}

func TestLogger(t *testing.T) {
	logger.Info("info")
}

func TestExec(t *testing.T) {
	output, err := Exec("go", "version")
	assert.NoError(t, err)
	t.Logf("%s", output)

	output, err = Exec("go", "mod", "init", "hello")
	assert.NoError(t, err)
	t.Logf("%s", output)
}
