package easylog

import (
	"testing"
	// "gotest.tools/assert"
)

func TestNewdefaultSamplerLogger(t *testing.T) {
	l := NewDefaultSamplerLogger()
	// assert.Equal(t, l, nil)
	l.Info("test", nil)
	l.Close()
}
