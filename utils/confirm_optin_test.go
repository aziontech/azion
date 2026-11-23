package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withStdin makes os.Stdin read input for the duration of the test.
func withStdin(t *testing.T, input string) {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	_, err = w.WriteString(input)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	old := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = old; _ = r.Close() })
}

func TestConfirmOptInUnattendedAnswersNo(t *testing.T) {
	// Nothing to read: a prompt would block or answer from stdin.
	withStdin(t, "y\n")
	assert.False(t, ConfirmOptIn(true, "question? (y/N)"))
}

func TestConfirmOptInDefaultsToNo(t *testing.T) {
	withStdin(t, "\n")
	assert.False(t, ConfirmOptIn(false, "question? (y/N)"))
}

func TestConfirmOptInAcceptsAnExplicitYes(t *testing.T) {
	withStdin(t, "y\n")
	assert.True(t, ConfirmOptIn(false, "question? (y/N)"))
}
