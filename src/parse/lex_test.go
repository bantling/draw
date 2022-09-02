package parse

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEof(t *testing.T) {
	src := strings.NewReader("")
	assert.Equal(t, cEof, Lex(src))
}
