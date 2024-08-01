package root

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestHasFailed(t *testing.T) {
	tests := []struct {
		name      string
		hasFailed bool
		expected  bool
	}{
		{
			name:      "HasFailed is true",
			hasFailed: true,
			expected:  true,
		},
		{
			name:      "HasFailed is false",
			hasFailed: false,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasFailed = tt.hasFailed

			result := HasFailed()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNestedSuggestFunc(t *testing.T) {
	tests := []struct {
		name           string
		arg            string
		setupCmd       func() *cobra.Command
		expectedOutput string
	}{
		{
			name: "arg is help",
			arg:  "help",
			setupCmd: func() *cobra.Command {
				return &cobra.Command{Use: "test"}
			},
			expectedOutput: "unknown command \"help\" for \"test\"\n\nDid you mean this?\n\t--help\n\n",
		},
		{
			name: "arg with suggestions",
			arg:  "hepl",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{Use: "test"}
				cmd.AddCommand(&cobra.Command{Use: "help"})
				return cmd
			},
			expectedOutput: "unknown command \"hepl\" for \"test\"\n\n",
		},
		{
			name: "arg without suggestions",
			arg:  "unknown",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{Use: "test"}
				return cmd
			},
			expectedOutput: "unknown command \"unknown\" for \"test\"\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			var output bytes.Buffer
			cmd.SetOut(&output)
			cmd.SetErr(&output)

			nestedSuggestFunc(cmd, tt.arg)

			assert.Equal(t, tt.expectedOutput, output.String())
		})
	}
}

func TestIsRootCmd(t *testing.T) {
	tests := []struct {
		name          string
		setupCmd      func() *cobra.Command
		expectedResult bool
	}{
		{
			name: "command is nil",
			setupCmd: func() *cobra.Command {
				return nil
			},
			expectedResult: false,
		},
		{
			name: "command is root",
			setupCmd: func() *cobra.Command {
				return &cobra.Command{Use: "root"}
			},
			expectedResult: true,
		},
		{
			name: "command has parent",
			setupCmd: func() *cobra.Command {
				parentCmd := &cobra.Command{Use: "root"}
				childCmd := &cobra.Command{Use: "child"}
				parentCmd.AddCommand(childCmd)
				return childCmd
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			result := isRootCmd(cmd)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
