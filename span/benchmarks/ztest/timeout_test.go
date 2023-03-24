package ztest

// This is a test file for the ztest package.
import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	base := 100 * time.Millisecond
	timeout := Timeout(base)
	assert.Equal(t, timeout, base, "Timeout should not be scaled.")
}

func TestSleep(t *testing.T) {
	start := time.Now()
	base := 100 * time.Millisecond
	Sleep(base)
	duration := time.Since(start)
	assert.True(t, duration >= base, "Sleep should block for at least the base duration.")
}

func TestInitialize(t *testing.T) {
	defer Initialize("1.0")()
	Initialize("2.0")()
	assert.Equal(t, Timeout(100*time.Millisecond), 100*time.Millisecond, "Timeout should be scaled.")
}
