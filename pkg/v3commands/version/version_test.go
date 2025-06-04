package version

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestVersionCommand(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name           string
		binVersion     string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "string",
			binVersion:     "development",
			expectedOutput: "Azion CLI development\n",
			expectError:    false,
		},
		{
			name:           "specific version",
			binVersion:     "1.2.3",
			expectedOutput: "Azion CLI 1.2.3\n",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the version
			BinVersion = tt.binVersion

			f, stdout, _ := testutils.NewFactory(nil)

			// Create the command
			versionCmd := NewCmd(f)

			// Execute the command
			_, err := versionCmd.ExecuteC()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, stdout.String())
			}
		})
	}
}
