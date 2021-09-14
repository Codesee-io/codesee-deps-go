package signals

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	ch := Setup()

	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	require.NoError(t, err)

	select {
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for signal")
	case <-ch:
	}
}
