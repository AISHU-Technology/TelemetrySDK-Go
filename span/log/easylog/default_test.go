package easylog

import (
    "testing"

    // "gotest.tools/assert"
)

func TestNewdefaultSamplerLogger(t *testing.T) {
    l := NewdefaultSamplerLogger()
    // assert.Equal(t, l, nil)
    l.Info("test", nil)
    l.Close()
}

