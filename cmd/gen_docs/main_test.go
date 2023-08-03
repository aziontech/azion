package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_run(t *testing.T) {
	t.Run("test creating dir and reading file without error", func(t *testing.T) {
		dir := t.TempDir()
		args := []string{"--doc-path", dir, "--file-type", "yaml"}
		err := run(args)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
	})

	t.Run("test no dir sent", func(t *testing.T) {
		args := []string{}
		err := run(args)
		require.Error(t, err, "error: --doc-path not set")
	})

	t.Run("test help cmd", func(t *testing.T) {
		args := []string{"--help"}
		err := run(args)
		require.NoError(t, err)
	})

}
