package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerator_Gen(t *testing.T) {
	g := NewGenerator("../../../go-starter-example")
	err := g.Gen()
	assert.NoError(t, err)
}
