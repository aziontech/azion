package cells

import (
	"errors"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/api/graphql/cells"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func mockCellsConsoleLogs(f *cmdutil.Factory, functionId string, logTime time.Time, limit string) (cells.CellsConsoleEventsResponse, error) {
	return cells.CellsConsoleEventsResponse{
		CellsConsoleEvents: []cells.CellsConsoleEvent{
			{
				FunctionId: "1234",
				Ts:         logTime,
				Level:      "LOG",
				Line:       "This is a log line",
			},
		},
	}, nil
}

func mockCellsConsoleLogsError(f *cmdutil.Factory, functionId string, logTime time.Time, limit string) (cells.CellsConsoleEventsResponse, error) {
	return cells.CellsConsoleEventsResponse{}, errors.New("error")
}

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}

	f, _, _ := testutils.NewFactory(mock)

	cmd := NewLogsCmd(f)
	cobracmd := NewCobraCmd(cmd)
	assert.NotNil(t, cmd)
	assert.IsType(t, &cobra.Command{}, cobracmd)
}

func TestPrintLogs(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name           string
		mockCells      func(*cmdutil.Factory, string, time.Time, string) (cells.CellsConsoleEventsResponse, error)
		expectedOutput string
		expectError    bool
		tail           bool
		pretty         bool
	}{
		{
			name:           "successful log retrieval",
			mockCells:      mockCellsConsoleLogs,
			expectedOutput: "Function ID: 1234",
			expectError:    false,
			tail:           false,
			pretty:         true,
		},
		{
			name:           "error in log retrieval",
			mockCells:      mockCellsConsoleLogsError,
			expectedOutput: "",
			expectError:    true,
			tail:           false,
			pretty:         true,
		},
		{
			name:           "pretty print logs",
			mockCells:      mockCellsConsoleLogs,
			expectedOutput: "Function ID: 1234",
			expectError:    false,
			tail:           false,
			pretty:         true,
		},
		{
			name:           "tail logs",
			mockCells:      mockCellsConsoleLogs,
			expectedOutput: "Function ID: 1234",
			expectError:    false,
			tail:           true,
			pretty:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var arguments []string

			mock := &httpmock.Registry{}

			f, stdout, _ := testutils.NewFactory(mock)

			cmd := NewLogsCmd(f)
			cmd.GetLogs = tt.mockCells
			cmd.Tail = tt.tail
			cobracmd := NewCobraCmd(cmd)
			if tt.pretty {
				arguments = append(arguments, "--pretty")
			}

			cobracmd.SetArgs(arguments)

			err := cobracmd.Execute()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				output := stdout.String()
				assert.Contains(t, output, tt.expectedOutput)
			}
		})
	}
}
