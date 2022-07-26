package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerator_Gen(t *testing.T) {
	g := NewGenerator("../../../go-starter-example", "github.com/wzzfarewell/go-starter-example")
	err := g.Gen("../testdata/config.toml")
	assert.NoError(t, err)
}
