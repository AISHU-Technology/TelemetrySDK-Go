package easylog

import (
	"testing"

	"gotest.tools/assert"
)

func NewdefaultSamplerLoggerTest(t *testing.T) {
	l := NewdefaultSamplerLogger()
	assert.Equal(t, l, nil)

	l.Close()
}
